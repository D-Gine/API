/*
** D&GINE Project, 2026
** API
** File description:
** rulesets/register.go
 */

package rulesets

import (
	"database/sql"
	"net/http"

	"api/src/database"
	"api/src/internal/structs"

	"github.com/gin-gonic/gin"
)

type CreateRulesetsArgs struct {
	Name string `json:"name"`
}

// @BasePath /api/rulesets
// Rulesets godoc
// @Summary Creates a ruleset
// @Schemes
// @Description Creates a ruleset with arguments in body<br><br>Will return the id of created ruleset
// @Tags rulesets
// @Accept json
// @Produce json
// @Param creds body CreateRulesetsArgs true "ruleset related informations that will be later needed for the login process"
// @Success 201 {object} structs.PostResponse
// @Router /api/rulesets [post]
func CreateRulesets(c *gin.Context) {
	var args CreateRulesetsArgs

	// Body Json parsing
	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Checking db if ruleset already exists
	rows, err := database.Db.Query("SELECT name FROM games.rulesets WHERE name=$1", args.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()
	if rows.Next() {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Ruleset already exists"})
		return
	}

	// Inserting new ruleset into db
	var id sql.NullString
	err = database.Db.QueryRow("INSERT INTO games.rulesets (name) VALUES ($1) RETURNING id", args.Name).Scan(&id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register ruleset"})
		return
	}

	c.JSON(http.StatusCreated, structs.PostResponse{Id: id.String})
}
