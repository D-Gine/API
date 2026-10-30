/*
** D&GINE Project, 2026
** API
** File description:
** users/read.go
 */

package users

import (
	"api/src/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/users
// Users godoc
// @Summary Reads all users data
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Reads all users data
// @Tags users
// @Produce json
// @Success 200 {object} []UserReadResponse
// @Router /api/users [get]
func ReadUsers(c *gin.Context) {
	result := []UserReadResponse{}
	rows, err := database.Db.Query("SELECT id, name, email, image, role FROM accounts.users")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id, name, email, image, role sql.NullString
		if err := rows.Scan(&id, &name, &email, &image, &role); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		result = append(result, UserReadResponse{
			Id:    id.String,
			Name:  name.String,
			Email: email.String,
			Image: image.String,
			Role:  role.String,
		})
	}

	c.JSON(http.StatusOK, result)
}
