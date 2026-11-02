package components

import (
	"api/src/controllers/shared_types"
	"encoding/json"
)

type ValidateValue struct {
	Value      json.RawMessage           `json:"value" swaggertype:"obect"`
	Components shared_types.ComponentMap `json:"components"`
}

type ValidateValues struct {
	Values     shared_types.ComponentMap `json:"values" swaggertype:"obect"`
	Components shared_types.ComponentMap `json:"components"`
}
