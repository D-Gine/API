package shared_types

import "encoding/json"

type NodeList struct {
	NodeValues map[string]json.RawMessage `json:"node_values" swaggertype:"object"`
}

type MultipleNodeValidation struct {
	NodeList
	NodeIds []string `json:"node_ids"`
}
