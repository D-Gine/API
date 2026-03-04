/*
** D&GINE Project, 2026
** Backend
** File description:
** user/update.go
 */

package users

import (
	"api/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateArgs struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Role  string `json:"role"`
}

// @BasePath /api/users
// Users godoc
// @Summary Updates given ID user with arguments in body
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Updates the given user's account data
// @Tags users
// @Accept json
// @Param id query string true "id of the user to do the action on"
// @Param data body UpdateArgs true "New arguments to be set in the user given, along with its id"
// @Success 200
// @Router /api/users [put]
func UpdateUsers(c *gin.Context) {
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

	_, err = database.Db.Exec(
		`UPDATE accounts.users
		SET email=$2, name=$3, image=$4, role=$5
		WHERE id=$1`,
		idQuery, args.Email, args.Name, args.Image, args.Role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
