package validate

import (
	"encoding/json"
)

type NodeList struct {
	NodeValues map[string]json.RawMessage `json:"node_values"`
}

type MultipleNodeValidation struct {
	NodeList
	NodeIds []string `json:"node_ids"`
}

type Verifications struct {
	NodeId       string          `json:"node_id"`
	BaseName     string          `json:"base_name"`
	Scripts      []string        `json:"scripts"`
	Restrictions json.RawMessage `json:"restrictions"`
}
