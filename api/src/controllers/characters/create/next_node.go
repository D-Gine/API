/*
** D&GINE Project, 2026
** API
** File description:
** characters/creation/next_node.go
 */

package charactersCreate

import (
	"api/src/database"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type ValuePair struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type NextNodeArgs struct {
	CurrentNodeId string      `json:"current_node_id"`
	Values        []ValuePair `json:"values"`
}

type ComponentTemplate struct {
	TemplateId string                 `json:"template_id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type NextNodeResponse struct {
	NodeId         string              `json:"node_id"`
	Name           string              `json:"name"`
	Last           bool                `json:"last"`
	ComponentTypes []ComponentTemplate `json:"component_types"`
}

type Condition struct {
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type PossibleValue struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func evaluateCondition(condition Condition, value interface{}) bool {
	if condition.Operator == "" {
		return true
	}

	valueStr, ok := value.(string)
	if !ok {
		return false
	}

	switch condition.Operator {
	case "==":
		return valueStr == condition.Value
	case "!=":
		return valueStr != condition.Value
	default:
		return false
	}
}

func findMatchingChild(currentNodeId string, values []ValuePair) (string, error) {
	rows, err := database.Db.Query(`
		SELECT DISTINCT child_id
		FROM games.creations_links
		WHERE parent_id = $1 AND child_id IS NOT NULL
	`, currentNodeId)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	var childIds []string
	for rows.Next() {
		var childId string
		if err := rows.Scan(&childId); err != nil {
			return "", err
		}
		childIds = append(childIds, childId)
	}

	if len(childIds) == 0 {
		return "", sql.ErrNoRows
	}

	for _, childId := range childIds {
		var conditionJson sql.NullString
		err := database.Db.QueryRow(`
			SELECT condition
			FROM games.components_creations
			WHERE id = $1
		`, childId).Scan(&conditionJson)

		if err != nil {
			continue
		}
		if !conditionJson.Valid || conditionJson.String == "{}" || conditionJson.String == "" {
			continue
		}
		var condition Condition
		if err := json.Unmarshal([]byte(conditionJson.String), &condition); err != nil {
			continue
		}
		var templateId string
		err = database.Db.QueryRow(`
			SELECT template_id
			FROM games.creations_templates
			WHERE creation_id = $1
			LIMIT 1
		`, currentNodeId).Scan(&templateId)

		if err != nil {
			continue
		}
		for _, valuePair := range values {
			if valuePair.Key == templateId {
				if evaluateCondition(condition, valuePair.Value) {
					return childId, nil
				}
				break
			}
		}
	}
	for _, childId := range childIds {
		var conditionJson sql.NullString
		err := database.Db.QueryRow(`
			SELECT condition
			FROM games.components_creations
			WHERE id = $1
		`, childId).Scan(&conditionJson)

		if err == nil && (!conditionJson.Valid || conditionJson.String == "{}" || conditionJson.String == "") {
			return childId, nil
		}
	}

	return "", sql.ErrNoRows
}

func getNodeTemplates(nodeId string) ([]ComponentTemplate, error) {
	rows, err := database.Db.Query(`
		SELECT
			ct.id,
			ct.name,
			ctype.type,
			ctype.metadata
		FROM games.creations_templates crt
		INNER JOIN games.components_templates ct ON crt.template_id = ct.id
		INNER JOIN games.components_types ctype ON ct.type = ctype.id
		WHERE crt.creation_id = $1
	`, nodeId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []ComponentTemplate
	for rows.Next() {
		var templateId, name, typeStr string
		var metadataJson sql.NullString

		if err := rows.Scan(&templateId, &name, &typeStr, &metadataJson); err != nil {
			return nil, err
		}

		template := ComponentTemplate{
			TemplateId: templateId,
			Name:       name,
			Type:       typeStr,
			Metadata:   make(map[string]interface{}),
		}

		if metadataJson.Valid && metadataJson.String != "" {
			if err := json.Unmarshal([]byte(metadataJson.String), &template.Metadata); err == nil {
				if typeStr == "enum_tag" {
					if tags, ok := template.Metadata["tags"].([]interface{}); ok {
						possibleValues, err := getPossibleValuesByTags(tags)
						if err == nil {
							template.Metadata["possible_values"] = possibleValues
						} else {
							template.Metadata["possible_values"] = []PossibleValue{}
						}
					}
				}
			}
		}

		templates = append(templates, template)
	}

	return templates, nil
}

func getPossibleValuesByTags(tags []interface{}) ([]PossibleValue, error) {
	if len(tags) == 0 {
		return []PossibleValue{}, nil
	}
	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tagStr, ok := tag.(string); ok {
			tagNames = append(tagNames, tagStr)
		}
	}

	if len(tagNames) == 0 {
		return []PossibleValue{}, nil
	}
	query := `
		SELECT DISTINCT i.id, i.name
		FROM games.items i
		INNER JOIN games.items_tags it ON i.id = it.item_id
		INNER JOIN games.tags t ON it.tag_id = t.id
		WHERE t.name = ANY($1)
	`

	rows, err := database.Db.Query(query, pq.Array(tagNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var possibleValues []PossibleValue
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			continue
		}
		possibleValues = append(possibleValues, PossibleValue{
			Id:   id,
			Name: name,
		})
	}

	return possibleValues, nil
}

func hasChildren(nodeId string) (bool, error) {
	var count int
	err := database.Db.QueryRow(`
		SELECT COUNT(*)
		FROM games.creations_links
		WHERE parent_id = $1 AND child_id IS NOT NULL
	`, nodeId).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// @BasePath /api/characters/create/nextnode
// Characters godoc
// @Summary Returns the next node
// @Schemes
// @Description Returns the next node of the given node depending on the ruleset character creation tree<br><br>This will check for the tree for all children of the given node, and return the one that matches the condition with the given values<br><br>If no child matches the condition, it means that the character creation process is finished
// @Tags characters creation
// @Accept json
// @Produce json
// @Param creds body NextNodeArgs true "Current node id and values for templates"
// @Success 200 {object} NextNodeResponse
// @Router /api/characters/create/nextnode [post]
func CharacterNextNode(c *gin.Context) {
	var args NextNodeArgs

	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	nextNodeId, err := findMatchingChild(args.CurrentNodeId, args.Values)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, NextNodeResponse{
				NodeId:         args.CurrentNodeId,
				Name:           "",
				Last:           true,
				ComponentTypes: []ComponentTemplate{},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	templates, err := getNodeTemplates(nextNodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get node templates: " + err.Error()})
		return
	}
	hasChild, err := hasChildren(nextNodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check children: " + err.Error()})
		return
	}

	nodeName := ""
	if len(templates) > 0 {
		nodeName = "Node " + nextNodeId[:8]
	}

	response := NextNodeResponse{
		NodeId:         nextNodeId,
		Name:           nodeName,
		Last:           !hasChild,
		ComponentTypes: templates,
	}

	c.JSON(http.StatusOK, response)
}
