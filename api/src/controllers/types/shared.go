package types

import (
	"api/src/controllers/shared_types"
	"encoding/json"
)

type ValidateValue struct {
	Value      json.RawMessage           `json:"value" swaggertype:"object"`
	Components shared_types.ComponentMap `json:"components"`
}

type ValidateValues struct {
	Values     shared_types.TypeValueMap `json:"values" swaggertype:"object"`
	Components shared_types.ComponentMap `json:"components"`
}
