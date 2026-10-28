package nodes

import (
	"api/src/database"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

const maxTypeDepth = 5

// completes a type so the client does not need extra calls:
//   - enum:       options = every enum of the ruleset carrying all the required tags
//   - list:       item    = the fully resolved type of the elements
func resolveTypeDetails(t *NodeTypeResponse, rulesetID string, depth int) error {
	if depth > maxTypeDepth {
		return fmt.Errorf("type nesting too deep")
	}

	switch t.Base {
	case "enum":
		tags := []string{}
		if raw, ok := t.Restrictions["tags"].([]any); ok {
			for _, tag := range raw {
				if s, ok := tag.(string); ok {
					tags = append(tags, s)
				}
			}
		}
		options, err := loadEnumOptions(rulesetID, tags)
		if err != nil {
			return err
		}
		t.Options = options

	case "list":
		itemTypeID, ok := t.Restrictions["type"].(string)
		if !ok || itemTypeID == "" {
			return nil
		}
		item, err := loadTypeByID(itemTypeID)
		if err != nil {
			return err
		}
		if err := resolveTypeDetails(&item, rulesetID, depth+1); err != nil {
			return err
		}
		t.Item = &item
	}
	return nil
}

func loadTypeByID(typeID string) (NodeTypeResponse, error) {
	var out NodeTypeResponse
	var raw []byte
	err := database.Db.QueryRow(
		`SELECT COALESCE(get_base_type_json($1::uuid), '{}'::jsonb)::text`, typeID,
	).Scan(&raw)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, fmt.Errorf("invalid type payload from database for type %s", typeID)
	}
	return out, nil
}

func loadEnumOptions(rulesetID string, tags []string) ([]EnumOption, error) {
	rows, err := database.Db.Query(`
		SELECT e.id::text, COALESCE(e.name, ''), COALESCE(e.description, '')
		FROM public.enums e
		WHERE e.ruleset_id = $1
		  AND (
			SELECT COUNT(DISTINCT l.tag_id)
			FROM public.link_enums_tags l
			WHERE l.enum_id = e.id AND l.tag_id = ANY($2::uuid[])
		  ) = cardinality($2::uuid[])
		ORDER BY e.name
	`, rulesetID, pq.Array(tags))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := []EnumOption{}
	for rows.Next() {
		var o EnumOption
		if err := rows.Scan(&o.Id, &o.Name, &o.Description); err != nil {
			return nil, err
		}
		options = append(options, o)
	}
	return options, rows.Err()
}
