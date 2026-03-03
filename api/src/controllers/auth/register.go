/*
** D&GINE Project, 2026
** Backend
** File description:
** auth/register.go
 */

package auth

import (
	"database/sql"
	"net/http"

	"api/src/database"
	"api/src/internal/structs"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterArgs struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Param email body string true "Email of the dedicated user"
// Param password body string true "Password of the corresponding user"

// @BasePath /api/auth/register
// Auth godoc
// @Summary Creation and connection of a new user
// @Schemes
// @Description Creation and connection of a new user <br>Token will be set in cookies if the arguments combination is valid
// @Tags auth
// @Accept json
// @Produce json
// @Param creds body RegisterArgs true "User related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostLoginResponse
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

	// Inserting new user into db
	var id, role sql.NullString
	err = database.Db.QueryRow("INSERT INTO accounts.users (email, name, password) VALUES ($1, $2, $3) RETURNING id, role", args.Email, args.Username, string(hashedPassword)).Scan(&id, &role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	var result structs.PostLoginResponse
	result.Id = id.String
	result.Token, err = BuildToken(c, UserData{id.String, args.Email, args.Username, role.String})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
