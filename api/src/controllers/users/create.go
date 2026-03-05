/*
** D&GINE Project, 2026
** API
** File description:
** users/register.go
 */

package users

import (
	"database/sql"
	"net/http"

	"api/src/database"
	"api/src/internal/structs"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CreateUsersArgs struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// @BasePath /api/users
// Users godoc
// @Summary Creates a user
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Creates a user with arguments in body
// @Tags users
// @Accept json
// @Produce json
// @Param creds body CreateUsersArgs true "User related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router /api/users [post]
func CreateUsers(c *gin.Context) {
	var args CreateUsersArgs

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
	var id sql.NullString
	err = database.Db.QueryRow("INSERT INTO accounts.users (email, name, password, role) VALUES ($1, $2, $3, $4) RETURNING id", args.Email, args.Username, string(hashedPassword), args.Role).Scan(&id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, structs.PostResponse{Id: id.String})
}
