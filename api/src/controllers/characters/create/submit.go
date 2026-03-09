/*
** D&GINE Project, 2026
** API
** File description:
** auth/register.go
 */

package charactersCreate

import (
	"api/src/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateArgs struct {
	RulesetId string `json:"ruleset_id"`
}

// @BasePath /api/characters/create/submit
// Characters godoc
// @Summary Submits a character list of components
// @Schemes
// @Description Submits a character list of components and creates the character in database<br><br>Will return the id of created character
// @Tags characters creation
// @Accept json
// @Produce json
// @Param creds body CreateArgs true "character related informations that will be later needed for the login process"
// @Success 201 {object} RulesetFirstNode
// @Router /api/characters/create/submit [post]
func SubmitCharacter(c *gin.Context) {
	var args CreateArgs

	// Body Json parsing
	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	c.JSON(http.StatusCreated, structs.PostResponse{
		Id: "",
	})
}
