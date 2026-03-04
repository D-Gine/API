/*
** D&GINE Project, 2026
** Backend
** File description:
** internal/structs/structs.go
 */

package auth

type PostLoginResponse struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}
