package validate

import (
	"encoding/json"
)

type NodeList struct {
	NodeValues map[string]json.RawMessage `json:"node_values"`
}

type ComponentMap = map[string]json.RawMessage

type MultipleNodeValidation struct {
	NodeList
	NodeIds []string `json:"node_ids"`
}

type BaseVerifications struct {
	BaseName     string          `json:"base_name"`
	Scripts      []string        `json:"scripts"`
	Restrictions json.RawMessage `json:"restrictions"`
}

type ComponentVerification struct {
	ComponentId string `json:"node_id"`
	BaseVerifications
}

type TypeVerifications struct {
	TypeId string `json:"type_id"`
	BaseVerifications
}

type Verification interface {
	GetBase() BaseVerifications
}

func (n ComponentVerification) GetBase() BaseVerifications {
	return n.BaseVerifications
}

func (t TypeVerifications) GetBase() BaseVerifications {
	return t.BaseVerifications
}
