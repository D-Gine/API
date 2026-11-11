/*
** D&GINE Project, 2026
** API
** File description:
** characters/submit.go
 */

package characters

import (
	"api/src/controllers/shared_types"
	"api/src/database"
	"encoding/json"
	"maps"
	"net/http"
	"slices"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	typegeneration "github.com/D-Gine/Gestionnaire2Type/type_generation"
	"github.com/D-Gine/Gestionnaire2Type/validate"
	"github.com/gin-gonic/gin"
)

type GenerateDynamicInfos struct {
	StaticValues gdt.ComponentMap                          `json:"static_values"`
	ToGenerate   map[string]typegeneration.GenerationInfos `json:"to_generate"`
}

// @BasePath /api/:ruleset_id/character/:id/submit
// Nodes godoc
// @Summary Characters - Submit
// @Schemes
// @Description save a character to db and check if valid<br>save the static components into the database and call the type validation scripts on them.
// @Tags characters
// @Accept json
// @Produce json
// @Param creds body string true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]string "Empty Object"
// @Router /api/:ruleset_id/character/:id/submit [put]
func SubmitCharacter(c *gin.Context) {
	id := c.Param("id")
	var body validate.NodeList

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	error_list, err := validate.ValidateNodesTypes(database.Db, body.NodeValues, slices.Collect(maps.Keys(body.NodeValues)), c.Param("ruleset_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(error_list) != 0 {
		c.JSON(http.StatusUnauthorized, error_list)
		return
	}
	if code, err := WriteCharacterInDb(id, body); err != nil {
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
	dynamic_map := make(gdt.ComponentMap)
	for id, infos := range res.ToGenerate {
		value, err := typegeneration.GenerateKnownedType(database.Db, res.StaticValues, infos)
		if err != nil {
			c.JSON(shared_types.GenerationErrorCode(err), gin.H{"error": err.Error()})
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
