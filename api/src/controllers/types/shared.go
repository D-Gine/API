package types

import (
	"encoding/json"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

type ValidateValue struct {
	Value      json.RawMessage  `json:"value" swaggertype:"object"`
	Components gdt.ComponentMap `json:"components"`
}

type ValidateValues struct {
	Values     gdt.TypeValueMap `json:"values" swaggertype:"object"`
	Components gdt.ComponentMap `json:"components"`
}
