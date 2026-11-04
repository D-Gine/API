package typegeneration

import (
	"api/src/controllers/shared_types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BaseTypeGeneration(ctx *gin.Context, staticValues shared_types.ComponentMap, infos GenerationInfos) (json.RawMessage, int, error) {
	switch infos.TypeName {
	case "number":
		r, err := shared_types.ParseRestriction[shared_types.NumberRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return generateNumber(r)
	case "string":
		r, err := shared_types.ParseRestriction[shared_types.StringRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return generateString(r)
	case "enum":
		fmt.Println(infos.Restrictions.String())
		r, err := shared_types.ParseRestriction[shared_types.EnumRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		fmt.Println(r.Tags)
		return generateEnum(r)
	case "list":
		r, err := shared_types.ParseRestriction[shared_types.ListRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return generateList(r, ctx, staticValues)
	case "dict":
		r, err := shared_types.ParseRestriction[shared_types.DictionaryRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return generateDictionary(r, ctx, staticValues)
	}
	return nil, http.StatusBadRequest, fmt.Errorf("Unknow type %s", infos.TypeName)
}
