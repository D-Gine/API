package typegeneration

import (
	"api/src/controllers/shared_types"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"

	"github.com/lucasjones/reggen"
)

func generateString(restriction shared_types.StringRestrictions) (json.RawMessage, int, error) {
	if restriction.RegexExpression != nil && *restriction.RegexExpression != "" {
		generated, err := reggen.Generate(*restriction.RegexExpression, 10)
		if err != nil {
			return nil, http.StatusUnauthorized, fmt.Errorf("error generating regex: %w", err)
		}
		jsonString, err := json.Marshal(generated)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return jsonString, http.StatusOK, nil
	}

	minLen := 0
	if restriction.MinLength != nil {
		minLen = max(*restriction.MinLength, 0)
	}

	maxLen := minLen + (rand.Int() % 254) + 1
	if restriction.MaxLength != nil {
		maxLen = *restriction.MaxLength
	}
	if minLen > maxLen {
		minLen, maxLen = maxLen, minLen
	}
	span := int64(maxLen - minLen + 1)
	length := minLen + int(rand.Int64N(span))

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	sb.Grow(length)
	for range length {
		idx := rand.IntN(len(letterBytes))
		sb.WriteByte(letterBytes[idx])
	}

	jsonString, err := json.Marshal(sb.String())
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return jsonString, http.StatusOK, nil
}
