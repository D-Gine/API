/*
** D&GINE Project, 2026
** Backend
** File description:
** users/read_id.go
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
// @Summary Reads given user data
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Reads the given user data
// @Tags users
// @Param id query string true "id of the user to do the action on"
// @Produce json
// @Success 200 {object} UserReadResponse
// @Router /api/users/id [get]
func ReadUsersId(c *gin.Context) {
	idQuery := c.Request.URL.Query().Get("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}

	var id, name, email, image, role sql.NullString
	err := database.Db.QueryRow("SELECT id, name, email, image, role FROM accounts.users WHERE id=$1", idQuery).Scan(&id, &name, &email, &image, &role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error : " + err.Error()})
	}

	c.JSON(http.StatusOK, UserReadResponse{
		Id:    id.String,
		Name:  name.String,
		Email: email.String,
		Image: image.String,
		Role:  role.String,
	})

}
