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
			SELECT COUNT(*) FROM link_enums_tags WHERE enum_id = $1 AND tag_id = ANY($2::uuid[])`, value, pq.Array(*restrictions.Tags))
		var count int
		rows.Scan(&count)
		if count != len(*restrictions.Tags) {
			return errors.New("does not respect tag restrictions")
		}
	}
	return nil
}
