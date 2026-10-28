/*
** D&GINE Project, 2026
** API
** File description:
** nodes/shared.go
 */

package nodes

import "database/sql"

type EnumOption struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NodeTypeResponse struct {
	Id           string            `json:"id"`
	Base         string            `json:"base"`
	Restrictions map[string]any    `json:"restrictions"`
	Options      []EnumOption      `json:"options,omitempty"`
	Item         *NodeTypeResponse `json:"item,omitempty"`
}

type FirstNodeResponse struct {
	NodeId string           `json:"node_id"`
	Type   NodeTypeResponse `json:"type"`
}

type CurrentTree struct {
	NodeId string       `json:"node_id"`
	Next   *CurrentTree `json:"next,omitempty"`
}

type GetNextNodeRequest struct {
	NodeValues  map[string]any `json:"node_values"`
	CurrentTree *CurrentTree   `json:"current_tree"`
}

type NextNodeResponse struct {
	NodeId string            `json:"node_id"`
	Final  bool              `json:"final,omitempty"`
	Type   NodeTypeResponse  `json:"type"`
	Next   *NextNodeResponse `json:"next,omitempty"`
}

type GetNextNodeResponse struct {
	Next *NextNodeResponse `json:"next"`
}

type DbNodeData struct {
	NodeID       string
	NextScriptID sql.NullString
	Script       sql.NullString
	Params       map[string]any
	Type         NodeTypeResponse
}
