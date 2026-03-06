/*
** D&GINE Project, 2026
** API
** File description:
** account/delete.go
 */

package auth

import (
	"api/src/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAdmin(userid string) bool {
	var role sql.NullString

	err := database.Db.QueryRow("SELECT role FROM accounts.users WHERE id=$1", userid).Scan(&role)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	return (role.Valid && role.String == "admin")
}

func AdminMiddleware(c *gin.Context) {
	// Recuperation du user depuis le token
	user, err := GetUserFromToken(c)
	if err != nil {
		return
	}

	if !IsAdmin(user.Id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not an admin"})
		c.Abort()
		return
	}
	c.Next()
}
