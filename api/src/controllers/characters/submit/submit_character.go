package submit

import (
	"api/src/controllers/characters/save"
	"api/src/controllers/nodes"
	"api/src/controllers/shared_types"
	typegeneration "api/src/controllers/type_generation"
	"api/src/database"
	"encoding/json"
	"maps"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type GenerateDynamicInfos struct {
	StaticValues shared_types.ComponentMap                 `json:"static_values"`
	ToGenerate   map[string]typegeneration.GenerationInfos `json:"to_generate"`
}

// @BasePath /api/:ruleset_id/character/:id/submit
// Nodes godoc
// @Summary save a character to db and check if valid
// @Schemes
// @Description save the static components into the database and call the type validation scripts on them.
// @Tags character
// @Accept json
// @Produce json
// @Param creds body string true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]string "Empty Object"
// @Router /api/:ruleset_id/character/:id/submit [put]
func SumbitCharacter(c *gin.Context) {
	id := c.Param("id")
	var body shared_types.NodeList

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	err, error_list := nodes.ValidateNodesTypes(c, &body, slices.Collect(maps.Keys(body.NodeValues)), c.Param("ruleset_id"))
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
	dynamic_map := make(shared_types.ComponentMap)
	for id, infos := range res.ToGenerate {
		value, code, err := typegeneration.GenerateKnownedType(c, res.StaticValues, infos)
		if err != nil {
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}
		dynamic_map[id] = value
	}
	jsonComponents, err := json.Marshal(dynamic_map)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = database.Db.Exec(`SELECT update_dynamic_components($1, $2)`, id, jsonComponents)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}
