package validate

import (
	"api/src/database"
	"errors"

	"github.com/lib/pq"
)

type enumRestrictions struct {
	Tags *[]string `json:"tags,omitempty"`
}

func VerifyEnum(componentList *ComponentMap, value string, restrictions enumRestrictions) error {
	if restrictions.Tags != nil {
		rows := database.Db.QueryRow(`
			SELECT COUNT(*) FROM link_enums_tags WHERE enum_id = $1 AND tag_id = ANY($2::uuid[])`, value, pq.Array(*restrictions.Tags))
		var count int
		rows.Scan(&count)
		if count != len(*restrictions.Tags) {
			return errors.New("does not respect tag restrictions")
		}
	}
	return nil
}
