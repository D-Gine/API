/*
** D&GINE Project, 2026
** API
** File description:
** internal/domains/domains.go
 */

package config

import (
	"api/src/docs"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

var TokenDomain string
var Trustedproxies = []string{}
var AllowedOrigins = map[string]bool{}

func LoadAllowedOrigins() bool {
	rawOrigins := os.Getenv("ALLOWED_ORIGINS")

	for _, origin := range strings.Split(rawOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			AllowedOrigins[origin] = true
		}
	}
	return true
}

func LoadTrustedProxies() bool {
	rawProxies := os.Getenv("TRUSTED_PROXIES")

	for _, proxy := range strings.Split(rawProxies, ",") {
		proxy = strings.TrimSpace(proxy)
		if proxy != "" {
			Trustedproxies = append(Trustedproxies, proxy)
		}
	}
	return true
}

func LoadDocValues() bool {
	docs.SwaggerInfo.Version = "API v" + os.Getenv("VERSION")
	docs.SwaggerInfo.Host = os.Getenv("API_HOST") + ":" + os.Getenv("API_PORT")

	// Description
	docs.SwaggerInfo.Description = "**D&Gine** API routes documentation<br>"
	docs.SwaggerInfo.Description += "Please ask us for any inquiries : [@Resh](https://github.com/BaptiD) [@abelabel](https://github.com/abeldaverio) [@eto NaN](https://github.com/eto-NaN)<br>"

	docs.SwaggerInfo.Description += "<br><br><br>Golang version **" + runtime.Version()[2:] + "**"
	docs.SwaggerInfo.Description += "<br>Gin version **" + gin.Version[1:] + "**"

	return true
}

func LoadConfig() {
	LoadAllowedOrigins()
	LoadTrustedProxies()
	LoadDocValues()
	TokenDomain = os.Getenv("TOKEN_DOMAIN")
}
