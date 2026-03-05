/*
** D&GINE Project, 2026
** API
** File description:
** user/shared.go
 */

package characters

type RowReadResponse struct {
	Id       string `json:"id"`
	PlayerId string `json:"player_id"`
	Name     string `json:"name"`
}
