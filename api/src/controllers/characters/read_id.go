/*
** D&GINE Project, 2026
** API
** File description:
** characters/read_id.go
 */

package characters

import (
	"api/src/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ComponentValue struct {
	Key      string                 `json:"key"`      // template_id
	Type     string                 `json:"type"`     // type enum (number, string, enum_tag, etc.)
	Metadata map[string]interface{} `json:"metadata"` // type metadata
	Value    interface{}            `json:"value"`    // actual value from component
}

type Component struct {
	ComponentId string           `json:"component_id"`
	Name        string           `json:"name"`
	Values      []ComponentValue `json:"values"`
}

type CharacterReadResponse struct {
	Name       string      `json:"name"`
	Components []Component `json:"components"`
}

// @BasePath /api/characters
// Characters godoc
// @Summary Reads given character data with all components
// @Schemes
// @Description Reads the given character/entity data including all components and their values
// @Tags characters
// @Param id path string true "ID of the character to read"
// @Produce json
// @Success 200 {object} CharacterReadResponse
// @Router /api/characters/{id} [get]
func ReadCharactersId(c *gin.Context) {
	characterId := c.Param("id")
	if characterId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Character ID is required"})
		return
	}

	var entityName string
	err := database.Db.QueryRow(`
		SELECT name
		FROM games.entities
		WHERE id = $1
	`, characterId).Scan(&entityName)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	components, err := getEntityComponents(characterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get components: " + err.Error()})
		return
	}

	response := CharacterReadResponse{
		Name:       entityName,
		Components: components,
	}

	c.JSON(http.StatusOK, response)
}

func getEntityComponents(entityId string) ([]Component, error) {
	rows, err := database.Db.Query(`
		SELECT DISTINCT
			c.id as component_id,
			c.name as component_name,
			c.value as component_value
		FROM games.components c
		INNER JOIN games.components_entities ce ON c.id = ce.component_id
		WHERE ce.entity_id = $1
	`, entityId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []Component
	for rows.Next() {
		var componentId, componentName string
		var componentValueJson sql.NullString

		if err := rows.Scan(&componentId, &componentName, &componentValueJson); err != nil {
			continue
		}
		var componentValueMap map[string]interface{}
		if componentValueJson.Valid && componentValueJson.String != "" {
			if err := json.Unmarshal([]byte(componentValueJson.String), &componentValueMap); err != nil {
				fmt.Printf("Error unmarshaling component value: %v\n", err)
				continue
			}
		}

		values, err := getComponentTemplateValues(componentId, componentValueMap)
		if err != nil {
			continue
		}

		component := Component{
			ComponentId: componentId,
			Name:        componentName,
			Values:      values,
		}

		components = append(components, component)
	}

	return components, nil
}

func getComponentTemplateValues(componentId string, componentValueMap map[string]interface{}) ([]ComponentValue, error) {
	rows, err := database.Db.Query(`
		SELECT
			ct.id as template_id,
			ctype.type as type_enum,
			ctype.metadata as type_metadata
		FROM games.components_templates_links ctl
		INNER JOIN games.components_templates ct ON ctl.template_id = ct.id
		INNER JOIN games.components_types ctype ON ct.type = ctype.id
		WHERE ctl.component_id = $1
	`, componentId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []ComponentValue

	for rows.Next() {
		var templateId, typeEnum string
		var typeMetadataJson sql.NullString

		if err := rows.Scan(&templateId, &typeEnum, &typeMetadataJson); err != nil {
			continue
		}
		metadata := make(map[string]interface{})
		if typeMetadataJson.Valid && typeMetadataJson.String != "" {
			json.Unmarshal([]byte(typeMetadataJson.String), &metadata)
		}
		var value interface{}
		if componentValueMap != nil {
			if val, ok := componentValueMap[templateId]; ok {
				value = val
				fmt.Printf("Found value for template %s: %v\n", templateId, value)
			} else {
				fmt.Printf("No value found for template %s in map: %+v\n", templateId, componentValueMap)
			}
		}

		componentValue := ComponentValue{
			Key:      templateId,
			Type:     typeEnum,
			Metadata: metadata,
			Value:    value,
		}

		values = append(values, componentValue)
	}

	return values, nil
}
