package shared_types

import "encoding/json"

type ComponentMap = map[string]json.RawMessage

type ComponentVerification struct {
	ComponentId string `json:"component_id"`
	BaseVerifications
}

func (n ComponentVerification) GetBase() BaseVerifications {
	return n.BaseVerifications
}
