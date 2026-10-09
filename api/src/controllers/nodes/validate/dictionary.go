package validate

import (
	"api/src/database"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/lib/pq"
)

type dictionaryRestrictions struct {
	Interface *map[string]string `json:"interface,omitempty"`
}

func VerifyDictionary(componentList *ComponentMap, value map[string]json.RawMessage, restrictions dictionaryRestrictions) error {
	if restrictions.Interface != nil {
		// getting verifications data
		type_set := make(map[string]struct{})
		for _, type_id := range *restrictions.Interface {
			type_set[type_id] = struct{}{}
		}
		types := slices.Collect(maps.Keys(type_set))
		rows, err := database.Db.Query(`
			SELECT * FROM get_verifications($1)`, pq.Array(types))
		if err != nil {
			return fmt.Errorf("error fetching interfaces types, please recheck validity of you interface, %w", err)
		}
		defer rows.Close()
		verificationMap := make(map[string]TypeVerifications)
		for rows.Next() {
			var verification TypeVerifications
			if err := rows.Scan(
				&verification.TypeId,
				&verification.BaseName,
				&verification.Restrictions,
				pq.Array(&verification.Scripts),
			); err != nil {
				return fmt.Errorf("error scaning db row, %w", err)
			}
			verificationMap[verification.TypeId] = verification
		}
		if err = rows.Err(); err != nil {
			return fmt.Errorf("error reading db rows, %w", err)
		}
		// checking values
		for key, elem := range value {
			typeId, ok := (*restrictions.Interface)[key]
			if !ok {
				return fmt.Errorf(`unknow key "%s"`, key)
			}
			err := VerifyType(componentList, elem, verificationMap[typeId])
			if err != nil {
				return fmt.Errorf(`error at key "%s", %w`, key, err)
			}
			delete(*restrictions.Interface, key)
		}
		if len(*restrictions.Interface) > 0 {
			return fmt.Errorf(
				"missing keys: %s",
				strings.Join(slices.Collect(maps.Keys(*restrictions.Interface)), ", "),
			)
		}

	}
	return nil
}
