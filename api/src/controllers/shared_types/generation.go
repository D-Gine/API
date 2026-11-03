package shared_types

type GenerationMultipleBody struct {
	Ids        []string     `json:"ids"`
	Components ComponentMap `json:"components"`
}
