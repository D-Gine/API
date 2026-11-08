package components

import (
	"encoding/json"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

type ValidateValue struct {
	Value      json.RawMessage  `json:"value" swaggertype:"object"`
	Components gdt.ComponentMap `json:"components"`
}

type ValidateValues struct {
	Values     gdt.ComponentMap `json:"values" swaggertype:"object"`
	Components gdt.ComponentMap `json:"components"`
}
