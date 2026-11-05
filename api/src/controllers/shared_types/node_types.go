package shared_types

import "encoding/json"

type NodeMap = map[string]json.RawMessage

type NodeList struct {
	NodeValues NodeMap `json:"node_values" swaggertype:"object"`
}

type MultipleNodeValidation struct {
	NodeList
	NodeIds []string `json:"node_ids"`
}
