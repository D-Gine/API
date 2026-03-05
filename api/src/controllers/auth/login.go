/*
** D&GINE Project, 2026
** API
** File description:
** auth/login.go
 */

package auth

import (
	"database/sql"
	"net/http"

	"api/src/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @BasePath /api/auth/login
// Auth godoc
// @Summary Connection of a user to an existing account
// @Schemes
// @Description Connection of a user to an existing account<br>Token will be set in cookies if the arguments combination is valid
// @Tags auth
// @Accept json
// @Produce json
// @Param creds body Creds true "Email + Password combination of the corresponding user"
// @Success 200 {object} PostLoginResponse
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var creds Creds

	// Body Json parsing
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Checking db if email + password combination matches
	var id, email, username, image, role, hashedPassword sql.NullString
	err = database.Db.QueryRow("SELECT id, email, name, image, role, password FROM accounts.users WHERE email=$1",
		creds.Email).Scan(&id, &email, &username, &image, &role, &hashedPassword)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword.String), []byte(creds.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generating token
	var result PostLoginResponse
	result.Id = id.String
	result.Token, err = BuildToken(c, UserData{id.String, email.String, username.String, image.String, role.String})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
