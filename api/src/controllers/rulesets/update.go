/*
** D&GINE Project, 2026
** API
** File description:
** rulesets/update.go
 */

package rulesets

import (
	"api/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateArgs struct {
	Name string `json:"name"`
}

// @BasePath /api/rulesets
// rulesets godoc
// @Summary Updates given ruleset
// @Schemes
// @Description Updates given ruleset with arguments in body
// @Tags rulesets
// @Accept json
// @Param id query string true "id of the ruleset to update"
// @Param data body UpdateArgs true "New arguments to be set in the ruleset given"
// @Success 200
// @Router /api/rulesets/id [put]
func UpdateRulesets(c *gin.Context) {
	idQuery := c.Request.URL.Query().Get("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}

	var args UpdateArgs

	// Parsing des args
	err := c.ShouldBindJSON(&args)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err = database.Db.Exec("UPDATE games.rulesets SET name=$2 WHERE id=$1", idQuery, args.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
