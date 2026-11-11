package shared_types

import (
	"errors"
	"net/http"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

type GenerationMultipleBody struct {
	Ids        []string         `json:"ids"`
	Components gdt.ComponentMap `json:"components"`
}

func GenerationErrorCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if errors.Is(err, gdt.ErrDatabase) {
		return http.StatusInternalServerError
	}
	return http.StatusBadRequest
}
