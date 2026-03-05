/*
** D&GINE Project, 2026
** Backend
** File description:
** users/shared.go
 */

package users

type UserReadResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
	Role  string `json:"role"`
}
