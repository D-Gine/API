/*
** D&GINE Project, 2026
** API
** File description:
** auth/register.go
 */

package auth

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"api/src/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterArgs struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Image    string `json:"image"`
}

// @BasePath /api/auth/register
// Auth godoc
// @Summary Connection of a new user to a new account
// @Schemes
// @Description Connection of a new user to a new account<br>Token will be set in cookies if the username or email is not already taken<br><br>The image must be base64 encoded and must contain the file type in the header (ex: data:image/png;base64,...)
// @Tags auth
// @Accept json
// @Produce json
// @Param creds body RegisterArgs true "User related informations that will be later needed for the login process"
// @Success 201 {object} PostLoginResponse
// @Router /api/auth/register [post]
func Register(c *gin.Context) {
	var args RegisterArgs

	// Body Json parsing
	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Checking db if user already exists
	rows, err := database.Db.Query("SELECT email FROM accounts.users WHERE email=$1", args.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()
	if rows.Next() {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hashing password to insert into db
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	imgData := []byte{}
	imgExt := ""
	imgPath := ""
	if args.Image != "" {
		imgBase64 := args.Image

		if idx := strings.Index(imgBase64, ","); idx != -1 {
			prefix := imgBase64[:idx]     // "data:image/jpeg;base64"
			imgBase64 = imgBase64[idx+1:] // strip prefix from b64 payload

			// Extract extension from "data:image/<ext>;base64"
			if start := strings.Index(prefix, "/"); start != -1 {
				if end := strings.Index(prefix, ";"); end != -1 {
					imgExt = prefix[start+1 : end] // "jpeg", "png", "gif", ...
				}
			}
		}

		if idx := strings.Index(imgBase64, ","); idx != -1 {
			imgBase64 = imgBase64[idx+1:]
		}

		imgData, err = base64.StdEncoding.DecodeString(imgBase64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid image encoding"})
			return
		}
	}

	// Inserting new user into db
	var id, role sql.NullString
	err = database.Db.QueryRow(
		"INSERT INTO accounts.users (email, name, password) VALUES ($1, $2, $3) RETURNING id, role",
		args.Email, args.Username, string(hashedPassword)).Scan(&id, &role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	if args.Image != "" {
		if err := os.MkdirAll(os.Getenv("USERS_IMAGES_PATH"), os.ModePerm); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image directory"})
			return
		}

		imgPath = os.Getenv("USERS_IMAGES_PATH") + id.String + "." + imgExt

		if err := os.WriteFile(imgPath, imgData, 0644); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user image"})
			return
		}

		_, err = database.Db.Exec(
			"UPDATE accounts.users SET image=$1 WHERE id=$2", imgPath, id.String)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}
	}

	var result PostLoginResponse
	result.Id = id.String

	result.Token, err = BuildToken(c, UserData{id.String, args.Email, args.Username, imgPath, role.String})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusCreated, result)
	}
}
