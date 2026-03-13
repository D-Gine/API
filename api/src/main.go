/*
** D&GINE Project, 2026
** API
** File description:
** main.go
 */

package main

import (
	"api/src/controllers/account"
	"api/src/controllers/auth"
	"api/src/controllers/characters"
	charactersCreate "api/src/controllers/characters/create"
	"api/src/controllers/dev"
	"api/src/controllers/fileserver"
	"api/src/controllers/rulesets"
	"api/src/controllers/users"
	"api/src/database"
	_ "api/src/docs"
	"api/src/internal/apidata"
	"api/src/internal/config"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CORSMiddleware(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if config.AllowedOrigins[origin] {
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

func setTrustedProxies(r *gin.Engine) {
	if len(config.Trustedproxies) > 0 {
		r.SetTrustedProxies(config.Trustedproxies)
	} else {
		r.SetTrustedProxies(nil)
	}
}

func main() {
	port := ":" + os.Getenv("API_PORT")
	config.LoadConfig()
	database.ConnectDatabase()

	r := gin.New()
	r.Use(gin.Logger())   // Logs all requests
	r.Use(gin.Recovery()) // Catches any panics and returns a 500 error instead of crashing
	r.Use(CORSMiddleware)
	setTrustedProxies(r)

	{
		r.GET("/about.json", apidata.GetAbout)

		api := r.Group("/api", CORSMiddleware)
		api.GET("/health", apidata.GetHealthCheck)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", auth.Login)
			authGroup.POST("/register", auth.Register)
			authGroup.POST("/logout", auth.Logout)
		}

		devGroup := api.Group("/dev")
		{
			devGroup.POST("/callengine", dev.CallEngine)
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
			usersGroup.PUT("/id", users.UpdateUsers)
			usersGroup.DELETE("/id", users.DeleteUsers)
			usersGroup.GET("/id", users.ReadUsersId)
		}

		// charactersGroup := api.Group("/characters", auth.AuthenticateMiddleware, auth.AdminMiddleware)
		charactersGroup := api.Group("/characters")
		{
			charactersGroup.GET("", characters.ReadCharacters)
			charactersGroup.PUT("/id", characters.UpdateCharacters)
			charactersGroup.DELETE("/id", characters.DeleteCharacters)
			charactersGroup.GET("/:id", characters.ReadCharactersId)

			charactersCreationGroup := charactersGroup.Group("/create")
			{
				charactersCreationGroup.GET("/firstnode", charactersCreate.CharacterFirstNode)
				charactersCreationGroup.POST("/nextnode", charactersCreate.CharacterNextNode)
				charactersCreationGroup.POST("/submit", charactersCreate.SubmitCharacter)
			}
		}

		// rulesetsGroup := api.Group("/rulesets", auth.AuthenticateMiddleware, auth.AdminMiddleware)
		rulesetsGroup := api.Group("/rulesets")
		{
			rulesetsGroup.POST("", rulesets.CreateRulesets)
			rulesetsGroup.GET("", rulesets.ReadRulesets)
			rulesetsGroup.PUT("/:id", rulesets.UpdateRulesets)
			rulesetsGroup.DELETE("/:id", rulesets.DeleteRulesets)
			rulesetsGroup.GET("/:id", rulesets.ReadRulesetId)
		}

	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.DocExpansion("none"),
	))

	img := r.Group("/img", CORSMiddleware)
	img.GET("/*filepath", fileserver.GetImage)

	fmt.Println("Listening and serving HTTP on " + port)
	r.Run(port)
}
