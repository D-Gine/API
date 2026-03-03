/*
** D&GINE Project, 2026
** Backend
** File description:
** internal/domains/domains.go
 */

package domains

const (
	TokenDomain = ""
)

var Trustedproxies = []string{
	"127.0.0.1"}

var AllowedOrigins = map[string]bool{
	"http://localhost:8081": true,
	"http://127.0.0.1:8081": true}
