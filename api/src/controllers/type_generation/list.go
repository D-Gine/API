package typegeneration

import (
	"api/src/controllers/shared_types"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func generateList(restriction shared_types.ListRestrictions, ctx *gin.Context, staticValues shared_types.ComponentMap) (json.RawMessage, int, error) {
	// TODO
	return nil, http.StatusBadRequest, errors.New("TODO")
}
