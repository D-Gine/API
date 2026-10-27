package typegeneration

import (
	"api/src/controllers"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func generateDictionary(restriction controllers.DictionaryRestrictions, ctx *gin.Context, staticValues controllers.ComponentMap) (json.RawMessage, int, error) {
	// TODO
	return nil, http.StatusBadRequest, errors.New("TODO")
}
