//Apenas para eu saber quais são as funções globais do projeto

package global

import (
	"encoding/json"
	"fmt"
	"strings"
)

func QuoteTokens(search []string) []string {
	result := make([]string, len(search))

	for i, value := range search {
		tokens := strings.Fields(value)

		for j := range tokens {
			tokens[j] = `"` + tokens[j] + `"`
		}

		result[i] = strings.Join(tokens, " ")
	}

	return result
}

type JsonMap map[string]any

func (j *JsonMap) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}

	var data []byte

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("JSONMap: tipo inesperado %T", value)
	}

	return json.Unmarshal(data, j)
}
