/*
** D&GINE Project, 2026
** Backend
** File description:
** user/update.go
 */

package user

import (
	"api-web/src/controllers/auth"
	"api-web/src/database"
	"api-web/src/internal/structs"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateArgs struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

// @BasePath /api/user
// User godoc
// @Summary Updates the actual user
// @Schemes
// @Description Updates the actual user <br><br><b>⚠️ The user must be logged in<b>
// @Tags user
// @Accept json
// @Param id body UpdateArgs true "New arguments to be set in the user"
// @Success 200 {object} structs.PostLoginResponse
// @Router /api/user [put]
func UpdateUser(c *gin.Context) {
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

	var result structs.PostLoginResponse
	result.Id = id.String
	result.Token, err = auth.BuildToken(c, auth.UserData{Id: id.String, Email: args.Email, Name: args.Username, Role: role.String})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
