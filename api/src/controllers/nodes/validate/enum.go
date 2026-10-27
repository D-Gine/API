package validate

import (
	"api/src/controllers"
	"api/src/database"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func VerifyEnum(ctx *gin.Context, componentList *controllers.ComponentMap, value string, restrictions controllers.EnumRestrictions) error {
	if restrictions.Tags != nil {
		rows := database.Db.QueryRow(`
			SELECT COUNT(*) FROM link_enums_tags WHERE enum_id = $1 AND tag_id = ANY($2::uuid[])`, value, pq.Array(restrictions.Tags))
		var count int
<<<<<<< HEAD
		if err := rows.Scan(&count); err != nil {
			return err
		}
		if count != len(*restrictions.Tags) {
=======
		rows.Scan(&count)
		if count != len(restrictions.Tags) {
>>>>>>> 17a628e6f0efddbe36d3b5002720e19996b76ab3
			return errors.New("does not respect tag restrictions")
		}
	}
	return nil
}
