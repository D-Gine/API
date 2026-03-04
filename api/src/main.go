/*
** D&GINE Project, 2026
** Backend
** File description:
** main.go
 */

package main

import (
	"api/src/controllers/account"
	"api/src/controllers/auth"
	"api/src/controllers/users"
	"api/src/database"
	_ "api/src/docs"
	"api/src/internal/domains"
	"api/src/internal/endpoints"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CORSMiddleware(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if domains.AllowedOrigins[origin] {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	} else {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	}
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func main() {
	port := ":" + os.Getenv("API_PORT")

	r := gin.Default()
	database.ConnectDatabase()
	r.Use(CORSMiddleware)

	{
		r.GET("/about.json", endpoints.GetAbout)

		api := r.Group("/api", CORSMiddleware)
		api.GET("/health", endpoints.GetHealthCheck)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", auth.Login)
			authGroup.POST("/register", auth.Register)
			authGroup.POST("/logout", auth.Logout)
			authGroup.POST("/oauth", auth.OAuthLogin)
		}

		accountGroup := api.Group("/account", auth.AuthenticateMiddleware)
		{
			accountGroup.GET("", account.ReadAccount)
			accountGroup.PUT("", account.UpdateAccount)
			accountGroup.DELETE("", account.DeleteAccount)
		}

		usersGroup := api.Group("/users", auth.AuthenticateMiddleware, auth.AdminMiddleware)
		{
			usersGroup.POST("", users.CreateUsers)
			usersGroup.GET("", users.ReadUsers)
			usersGroup.PUT("", users.UpdateUsers)
			usersGroup.DELETE("", users.DeleteUsers)
			usersGroup.GET("/id", users.ReadUsersId)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))
	fmt.Println("Listening and serving HTTP on " + port)
	r.Run(port)
}
