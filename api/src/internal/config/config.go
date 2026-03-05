/*
** D&GINE Project, 2026
** Backend
** File description:
** internal/domains/domains.go
 */

package config

import (
	"os"
	"strings"
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

func LoadConfig() {
	LoadAllowedOrigins()
	LoadTrustedProxies()
	TokenDomain = os.Getenv("TOKEN_DOMAIN")
}
