package typegeneration

import (
	"api/src/controllers"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BaseTypeGeneration(ctx *gin.Context, staticValues controllers.ComponentMap, infos GenerationInfos) (json.RawMessage, int, error) {
	switch infos.TypeName {
	case "number":
		r, err := controllers.ParseRestriction[controllers.NumberRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return generateNumber(r)
	case "string":
		r, err := controllers.ParseRestriction[controllers.StringRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return generateString(r)
	case "enum":
		r, err := controllers.ParseRestriction[controllers.EnumRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return generateEnum(r)
	case "list":
		r, err := controllers.ParseRestriction[controllers.ListRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return generateList(r, ctx, staticValues)
	case "dict":
		r, err := controllers.ParseRestriction[controllers.DictionaryRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return generateDictionary(r, ctx, staticValues)
	}
	return nil, http.StatusBadRequest, fmt.Errorf("Unknow type %s", infos.TypeName)
}
