/*
** D&GINE Project, 2026
** API
** File description:
** user/update.go
 */

package account

import (
	"api/src/controllers/auth"
	"api/src/database"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

var allowedImageExts = map[string]string{
	"png":  "png",
	"jpg":  "jpg",
	"jpeg": "jpg",
	"webp": "webp",
}

// @BasePath /api/account
// Account godoc
// @Summary Account - Update
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch current DB values
	var dbEmail, dbName, dbImage sql.NullString
	err = database.Db.QueryRow(
		"SELECT email, name, image FROM accounts.users WHERE id=$1", user.Id,
	).Scan(&dbEmail, &dbName, &dbImage)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
		return
	}

	// Keep old values when the new ones are empty
	newEmail := args.Email
	if newEmail == "" {
		newEmail = dbEmail.String
	}
	newUsername := args.Username
	if newUsername == "" {
		newUsername = dbName.String
	}

	// Password (only hashed if provided)
	var newPasswordHash string
	if args.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("UpdateAccount: bcrypt failed: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		newPasswordHash = string(hashed)
	}

	// Image handling: the file is written first, the old one is deleted
	// only after the DB update succeeds.
	var (
		imageChanged bool
		newImagePath string
	)

	if args.Image != "" {
		imgBase64 := args.Image
		imgExt := ""

		if idx := strings.Index(imgBase64, ","); idx != -1 {
			prefix := imgBase64[:idx]
			imgBase64 = imgBase64[idx+1:]

			if start := strings.Index(prefix, "/"); start != -1 {
				if end := strings.Index(prefix, ";"); end != -1 && end > start {
					imgExt = strings.ToLower(prefix[start+1 : end])
				}
			}
		}

		normalizedExt, ok := allowedImageExts[imgExt]
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
			return
		}

		imgData, err := base64.StdEncoding.DecodeString(imgBase64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid image encoding"})
			return
		}

		imagesDir := os.Getenv("USERS_IMAGES_PATH")
		if err := os.MkdirAll(imagesDir, os.ModePerm); err != nil {
			log.Printf("UpdateAccount: mkdir failed: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image directory"})
			return
		}

		newImagePath = filepath.Join(imagesDir, user.Id+"."+normalizedExt)
		if err := os.WriteFile(newImagePath, imgData, 0600); err != nil {
			log.Printf("UpdateAccount: write image failed: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new image"})
			return
		}
		imageChanged = true
	}

	setClauses := []string{"email=$1", "name=$2"}
	params := []interface{}{newEmail, newUsername}

	if newPasswordHash != "" {
		params = append(params, newPasswordHash)
		setClauses = append(setClauses, fmt.Sprintf("password=$%d", len(params)))
	}
	if imageChanged {
		params = append(params, newImagePath)
		setClauses = append(setClauses, fmt.Sprintf("image=$%d", len(params)))
	}

	params = append(params, user.Id)
	query := fmt.Sprintf(
		"UPDATE accounts.users SET %s WHERE id=$%d RETURNING id, role, image",
		strings.Join(setClauses, ", "), len(params),
	)

	var id, role, image sql.NullString
	err = database.Db.QueryRow(query, params...).Scan(&id, &role, &image)
	if err != nil {
		log.Printf("UpdateAccount: update user %s failed: %v", user.Id, err)

		// Roll back the newly written file so we don't leave orphans
		if imageChanged && newImagePath != dbImage.String {
			_ = os.Remove(newImagePath)
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	// DB update succeeded; safe to update image
	if imageChanged && dbImage.String != "" && dbImage.String != newImagePath {
		if err := os.Remove(dbImage.String); err != nil && !os.IsNotExist(err) {
			log.Printf("UpdateAccount: failed to delete old image %s: %v", dbImage.String, err)
		}
	}

	token, err := auth.BuildToken(c, auth.UserData{
		Id:    id.String,
		Email: newEmail,
		Name:  newUsername,
		Image: image.String, // value actually stored in the DB
		Role:  role.String,
	})
	if err != nil {
		log.Printf("UpdateAccount: build token failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build token"})
		return
	}

	var result auth.PostLoginResponse
	result.Id = id.String
	result.Token = token

	c.JSON(http.StatusOK, result)
}
