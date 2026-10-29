/*
** D&GINE Project, 2026
** API
** File description:
** user/shared.go
 */

package characters

import (
	"api/src/controllers/nodes"
	"encoding/json"
)

type RowReadResponse struct {
	Id       string `json:"id"`
	PlayerId string `json:"player_id"`
	Name     string `json:"name"`
}

type Component struct {
	Name  string                 `json:"name"`
	Value json.RawMessage        `json:"value"`
	Type  nodes.NodeTypeResponse `json:"type"`
}

type ComponentMap = map[string]Component

type EntityComponents struct {
	StaticComponents  ComponentMap `json:"static_components"`
	DynamicComponents ComponentMap `json:"dynamic_components"`
}

type Entity struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Components  EntityComponents `json:"components"`
}
