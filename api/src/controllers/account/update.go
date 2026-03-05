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
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UpdateArgs struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Image    string `json:"image"`
}

// @BasePath /api/account
// Account godoc
// @Summary Updates the user's account
// @Schemes
// @Description <b>⚠️ The user must be logged in ⚠️</b><br><br>Updates the user's account<br><br>The image must be base64 encoded and must contain the file type in the header (ex: data:image/png;base64,...)<br><br>All fields are optionals, if a field is empty or not provided, it will not be updated
// @Tags account
// @Accept json
// @Param id body UpdateArgs true "New arguments to be set in the user"
// @Success 200 {object} auth.PostLoginResponse
// @Router /api/account [put]
func UpdateAccount(c *gin.Context) {
	var args UpdateArgs

	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := auth.GetUserFromToken(c)
	if err != nil {
		return
	}

	// Getting db infos
	var dbEmail, dbName, dbImage sql.NullString
	err = database.Db.QueryRow(
		"SELECT email, name, image FROM accounts.users WHERE id=$1", user.Id,
	).Scan(&dbEmail, &dbName, &dbImage)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
		return
	}

	// Merging new data with old data if new data is empty
	newEmail := args.Email
	if newEmail == "" {
		newEmail = dbEmail.String
	}
	newUsername := args.Username
	if newUsername == "" {
		newUsername = dbName.String
	}

	newPassword := ""
	if args.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		newPassword = string(hashed)
	}

	// Gestion de l'image
	newImagePath := dbImage.String
	if args.Image != "" {
		imgBase64 := args.Image
		imgExt := ""

		if idx := strings.Index(imgBase64, ","); idx != -1 {
			prefix := imgBase64[:idx]
			imgBase64 = imgBase64[idx+1:]

			if start := strings.Index(prefix, "/"); start != -1 {
				if end := strings.Index(prefix, ";"); end != -1 {
					imgExt = prefix[start+1 : end]
				}
			}
		}

		if imgExt == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
			return
		}

		imgData, err := base64.StdEncoding.DecodeString(imgBase64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid image encoding"})
			return
		}

		if dbImage.String != "" {
			if err := os.Remove(dbImage.String); err != nil && !os.IsNotExist(err) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old image"})
				return
			}
		}

		if err := os.MkdirAll(os.Getenv("USERS_IMAGES_PATH"), os.ModePerm); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image directory"})
			return
		}
		newImagePath = os.Getenv("USERS_IMAGES_PATH") + user.Id + "." + imgExt
		if err := os.WriteFile(newImagePath, imgData, 0644); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new image"})
			return
		}
	}

	var id, role sql.NullString
	if args.Password != "" {
		err = database.Db.QueryRow(
			"UPDATE accounts.users SET email=$1, name=$2, password=$3, image=$4 WHERE id=$5 RETURNING id, role",
			newEmail, newUsername, newPassword, newImagePath, user.Id,
		).Scan(&id, &role)
	} else {
		err = database.Db.QueryRow(
			"UPDATE accounts.users SET email=$1, name=$2, image=$3 WHERE id=$4 RETURNING id, role",
			newEmail, newUsername, newImagePath, user.Id,
		).Scan(&id, &role)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account", "details": err})
		return
	}

	var result auth.PostLoginResponse
	result.Id = id.String
	result.Token, err = auth.BuildToken(c, auth.UserData{
		Id:    id.String,
		Email: newEmail,
		Name:  newUsername,
		Image: newImagePath,
		Role:  role.String,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
