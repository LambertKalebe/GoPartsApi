//Apenas para eu saber quais são as funções globais do projeto

package global

import (
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

	fmt.Println("quoteTokens", result)

	return result
}
