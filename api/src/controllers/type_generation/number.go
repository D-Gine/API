package typegeneration

import (
	"api/src/controllers"
	"encoding/json"
	"math"
	"math/rand/v2"
	"net/http"
)

func generateNumber(restriction controllers.NumberRestrictions) (json.RawMessage, int, error) {
	min := int64(math.MinInt32)
	if restriction.Min != nil {
		min = int64(*restriction.Min)
	}

	max := int64(math.MaxInt32)
	if restriction.Max != nil {
		max = int64(*restriction.Max)
	}

	if min > max {
		min, max = max, min
	}
	valueRange := (max - min) + 1
	value := int32(min + rand.Int64N(valueRange))

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return jsonValue, http.StatusInternalServerError, err
	} else {
		return jsonValue, http.StatusOK, err
	}
}
