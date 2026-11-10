package rulesets

import (
	"api/src/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/:ruleset_id/spritesheet_config
// rulesets godoc
// @Summary Rulesets - get spritesheet config
// @Schemes
// @Description get spritesheet config with ruleset id
// @Tags rulesets
// @Produce json
// @Success 200 {object} map[string]any
// @Router /api/:ruleset_id/spritesheet_config [get]
func GetSpritesheet(c *gin.Context) {
	ruleset_id := c.Param("ruleset_id")

	var res json.RawMessage
	err := database.Db.QueryRow(`
		SELECT config FROM character_sheet WHERE ruleset_id = $1
	`, ruleset_id).Scan(&res)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
