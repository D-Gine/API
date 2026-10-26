package submit

import (
	"api/src/controllers/characters/save"
	nodeValidate "api/src/controllers/nodes/validate"
	"api/src/database"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type GenerateDynamicInfos struct {
	StaticValues nodeValidate.ComponentMap `json:"static_values"`
	ToGenerate   map[string]struct {
		Script       *string         `json:"script,omitempty"`
		TypeName     string          `json:"type_name"`
		Restrictions json.RawMessage `json:"restrictions"`
	} `json:"to_generate"`
}

// // @BasePath api/:ruleset_id/character/:id/submit
// // Nodes godoc
// // @Summary save a character to db
// // @Schemes
// // @Description save the static components into the database.
// // @Tags character
// // @Accept json
// // @Produce json
// // @Param creds body CreateRulesetsArgs true "ruleset related informations that will be later needed for the login process"
// // @Success 200 {object} structs.PostResponse
// // @Router api/:ruleset_id/character/:id/submit [put]
func SumbitCharacter(c *gin.Context) {
	id := c.Param("id")
	var body nodeValidate.NodeList

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	err, error_list := nodeValidate.ValidateNodesTypes(c, &body, slices.Collect(maps.Keys(body.NodeValues)), c.Param("ruleset_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(error_list) != 0 {
		c.JSON(http.StatusUnauthorized, error_list)
		return
	}
	if code, err := save.WriteCharacterInDb(id, body); err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	var resRaw json.RawMessage
	var res GenerateDynamicInfos
	err = database.Db.QueryRow(`
		SELECT get_dynamic_generation_infos($1)
	`, id).Scan(&resRaw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := json.Unmarshal(resRaw, &res); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(*res.ToGenerate["2bb88ea8-9fe1-4b7b-85ca-bc15aa7201d0"].Script)
}
