/*
** D&GINE Project, 2026
** Backend
** File description:
** user/update.go
 */

package account

import (
	"api/src/controllers/auth"
	"api/src/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateArgs struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

// @BasePath /api/account
// Account godoc
// @Summary Updates the user's account
// @Schemes
// @Description <b>⚠️ The user must be logged in ⚠️</b><br><br>Updates the user's account
// @Tags account
// @Accept json
// @Param id body UpdateArgs true "New arguments to be set in the user"
// @Success 200 {object} auth.PostLoginResponse
// @Router /api/account [put]
func UpdateAccount(c *gin.Context) {
	var args UpdateArgs

	// Parsing des args
	err := c.ShouldBindJSON(&args)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Recuperation du user depuis le token
	user, err := auth.GetUserFromToken(c)
	if err != nil {
		return
	}

	var id, role sql.NullString
	err = database.Db.QueryRow("UPDATE accounts.users SET email=$1, name=$2 WHERE email=$3 returning id, role", args.Email, args.Username, user.Email).Scan(&id, &role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email", "details": err})
		return
	}

	var result auth.PostLoginResponse
	result.Id = id.String
	result.Token, err = auth.BuildToken(c, auth.UserData{Id: id.String, Email: args.Email, Name: args.Username, Role: role.String})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
