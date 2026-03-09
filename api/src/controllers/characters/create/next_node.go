/*
** D&GINE Project, 2026
** API
** File description:
** characters/creation/next_node.go
 */

package charactersCreate

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NextNodeArgs struct {
	NodeId string `json:"node_id"`
}

// @BasePath /api/characters/create/nextnode
// Characters godoc
// @Summary Returns the next node
// @Schemes
// @Description Returns the next node of the given node depending on the ruleset character creation tree
// @Tags characters creation
// @Accept json
// @Produce json
// @Param creds body NextNodeArgs true "character related informations that will be later needed for the login process"
// @Success 201
// @Router /api/characters/create/nextnode [get]
func CharacterNextNode(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
