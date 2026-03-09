/*
** D&GINE Project, 2026
** API
** File description:
** rulesets/delete.go
 */

package rulesets

import (
	"api/src/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/rulesets
// rulesets godoc
// @Summary Deletes given ruleset
// @Schemes
// @Param id query string true "id of the ruleset to delete"
// @Description Deletes the given ruleset<br> <b>Careful, there is no turn back or check, once called, this route will delete the ruleset no matter what</b> <br> It also deletes all content related to the ruleset depending on the parameters of the database
// @Tags rulesets
// @Success 204
// @Router /api/rulesets/id [delete]
func DeleteRulesets(c *gin.Context) {
	idQuery := c.Request.URL.Query().Get("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}

	_, err := database.Db.Exec("DELETE FROM games.rulesets WHERE id=$1", idQuery)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error" + err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
