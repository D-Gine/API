/*
** D&GINE Project, 2026
** Backend
** File description:
** auth/token.go
 */

package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("TOKEN_KEY"))

type UserData struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Role  string `json:"role"`
}

func CreateToken(user UserData) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"email":    user.Email,
		"username": user.Name,
		"image":    user.Image,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})
	tokenString, err := claims.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ExtractFromToken(tokenString string) (UserData, error) {
	parsedToken, err := verifyToken(tokenString)
	if err != nil {
		return UserData{}, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return UserData{}, fmt.Errorf("invalid token claims")
	}

	id, ok := claims["id"].(string)
	if !ok {
		return UserData{}, fmt.Errorf("id not found in token claims")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return UserData{}, fmt.Errorf("email not found in token claims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return UserData{}, fmt.Errorf("username not found in token claims")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return UserData{}, fmt.Errorf("role not found in token claims")
	}
	image, ok := claims["image"].(string)
	if !ok {
		return UserData{}, fmt.Errorf("image not found in token claims")
	}
	return UserData{id, email, username, image, role}, nil
}

func AuthenticateMiddleware(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.Next()
		return
	}

	var tokenString string
	cookieToken, err := c.Cookie("token")
	if err == nil {
		tokenString = cookieToken
	} else {
		// on prend le token du header si pas de cookie
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("Token missing in cookie and header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login or Register to access this page"})
			c.Abort()
			return
		}
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			tokenString = authHeader
		}
	}

	_, err = verifyToken(tokenString)
	if err != nil {
		fmt.Println("Token verification failed")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Next()
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func BuildToken(c *gin.Context, userdata UserData) (string, error) {
	tokenString, err := CreateToken(userdata)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return "", err
	}

	c.SetCookie("token", tokenString, 3600, "/", "", false, true)
	return tokenString, nil
}

func GetUserFromToken(c *gin.Context) (UserData, error) {
	var tokenString string

	cookieToken, err := c.Cookie("token")
	if err == nil {
		tokenString = cookieToken
	} else {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			return UserData{}, fmt.Errorf("missing token")
		}

		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			tokenString = authHeader
		}
	}

	user, err := ExtractFromToken(tokenString)
	if err != nil {
		return UserData{}, fmt.Errorf("invalid token")
	}

	return user, nil
}
