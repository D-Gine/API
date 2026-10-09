/*
** D&GINE Project, 2026
** API
** File description:
** nodes/shared.go
 */

package nodes

type NodeTypeResponse struct {
	Id           string         `json:"id"`
	Base         string         `json:"base"`
	Restrictions map[string]any `json:"restrictions"`
}

type FirstNodeResponse struct {
	NodeId string           `json:"node_id"`
	Type   NodeTypeResponse `json:"type"`
}
