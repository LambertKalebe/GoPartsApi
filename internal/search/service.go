package search

import (
	"errors"
	"fmt"
	"strings"
)

// Problemas com a pesquisa atual:
// - Não retorna produtos com base em seus similares (Falta de implementação no FTS5 do banco)
// - Não possui contexto de tokens (diferenciar anos, carros, peças, etc)
// - Ordenação de resultados precisa ser melhorada
// - Não sabe reconhecer que KL582 e KL 582 é a mesma coisa, apenas tirar os espaços se torna impossivel
// pois faria com que pesquisas como "KL 582 Fox" virasse "KL582Fox",
// a solução provavelmente se baseara na inserção do codigo não normalizado no FTS5
func serviceSearchProducts(search string, limit int) (productSearchResponse, error) {

	if limit <= 0 {
		limit = 5
	} else if limit > 100 {
		limit = 100
	}

	search = strings.TrimSpace(search)
	r := strings.NewReplacer(" de ", " ", " da ", " ", " do ", " ")
	search = r.Replace(search)
	if search == "" {
		return productSearchResponse{}, errors.New("invalid query")
	}
	fmt.Println("service", search)
	rows, err := searchProducts(search, limit)
	if err != nil {
		fmt.Println("service", err)
		return productSearchResponse{}, err
	}
	res, err := toProductSearchResponse(rows)
	if err != nil {
		fmt.Println("service", err)
		return productSearchResponse{}, err
	}
	return res, nil
}
