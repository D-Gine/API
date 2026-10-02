package validate

import (
	"api/src/database"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

type listRestrictions struct {
	Type      string `json:"type"`
	MinLength *int   `json:"min_length,omitempty"`
	MaxLength *int   `json:"max_length,omitempty"`
}

func VerifyList(componentList *ComponentMap, value []json.RawMessage, restrictions listRestrictions) error {
	var verification TypeVerifications
	err := database.Db.QueryRow(`
		SELECT * FROM get_verifications($1)`, pq.Array([]string{restrictions.Type})).Scan(
		&verification.TypeId,
		&verification.BaseName,
		&verification.Restrictions,
		pq.Array(&verification.Scripts),
	)
	if err != nil {
		return fmt.Errorf("error fetching for base type, please check your restriction 'type' field, %w", err)
	}
	if restrictions.MaxLength != nil && len(value) > *restrictions.MaxLength {
		return fmt.Errorf("too long, max lenght is %d, current len is %d", *restrictions.MaxLength, len(value))
	}
	if restrictions.MinLength != nil && len(value) < *restrictions.MinLength {
		return fmt.Errorf("too short, min lenght is %d, current len is %d", *restrictions.MinLength, len(value))
	}
	for i, elem := range value {
		if err = VerifyType(componentList, elem, verification); err != nil {
			return fmt.Errorf("error at element : %d, %w", i, err)
		}
	}
	return nil
}
