package typegeneration

import (
	"api/src/controllers"
	"api/src/database"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lib/pq"
)

func generateEnum(restriction controllers.EnumRestrictions) (json.RawMessage, int, error) {
	var enum string

	if restriction.Tags == nil {
		//TODO
		return nil, http.StatusUnauthorized, errors.New("cannot generate enum with no tag")
	}
	err := database.Db.QueryRow(`
			SELECT enum_id
			FROM link_enums_tags
			WHERE tag_id = ANY($1::UUID[])
			GROUP BY enum_id
			HAVING COUNT(DISTINCT tag_id) = array_length($1, 1)
			ORDER BY RANDOM()
			LIMIT 1
		`, pq.Array(*restriction.Tags)).Scan(&enum)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusBadRequest, fmt.Errorf("No enum found for tag requirement %q", *restriction.Tags)
		}
	}
	return nil, 0, nil
}
