/*
** D&GINE Project, 2026
** API
** File description:
** nodes/shared.go
 */

package nodes

import (
	"api/src/controllers/shared_types"
	"database/sql"
)

type NodeTypeResponse struct {
	Id           string         `json:"id"`
	Base         string         `json:"base"`
	Restrictions map[string]any `json:"restrictions"`
}

// NodeInfos holds what describes a node, shared by every node response.
type NodeInfos struct {
	NodeId      string           `json:"node_id"`
	TemplateId  string           `json:"template_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        NodeTypeResponse `json:"type"`
}

type CurrentTree struct {
	NodeId string       `json:"node_id"`
	Next   *CurrentTree `json:"next,omitempty"`
}

type GetNextNodeRequest struct {
	shared_types.NodeList
	CurrentTree *CurrentTree `json:"current_tree"`
}

type NextNodeResponse struct {
	NodeInfos
	Final bool              `json:"final,omitempty"`
	Next  *NextNodeResponse `json:"next,omitempty"`
}

type GetNextNodeResponse struct {
	Next *NextNodeResponse `json:"next"`
}

type DbNodeData struct {
	NodeInfos
	NextScriptID sql.NullString
	Script       sql.NullString
	Params       map[string]any
}
