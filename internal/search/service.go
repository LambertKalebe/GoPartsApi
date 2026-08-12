package search

import (
	"errors"
	"fmt"
	"g0/internal/global"
	"strings"
)

// Problemas com a pesquisa atual:
// - Não retorna produtos com base em seus similares (Falta de implementação no FTS5 do banco)
// - Não possui contexto de tokens (diferenciar anos, carros, peças, etc)
// - Ordenação de resultados precisa ser melhorada
// - Não sabe reconhecer que KL582 e KL 582 é a mesma coisa, apenas tirar os espaços se torna impossivel
// pois faria com que pesquisas como "KL 582 Fox" virasse "KL582Fox",
// a solução provavelmente se baseara na inserção do código não normalizado no FTS5
func serviceSearchProducts(search string, limit int) (productSearchResponse, error) {

	if limit <= 0 {
		limit = 5
	} else if limit > 100 {
		limit = 100
	}

	search = strings.TrimSpace(search)
	r := strings.NewReplacer(
		" de ", " ",
		" da ", " ",
		" do ", " ")
	search = r.Replace(search)
	if search == "" {
		return productSearchResponse{}, errors.New("invalid query")
	}
	quoted := global.QuoteTokens([]string{search})[0]
	fmt.Println("service", quoted)
	rows, err := searchProducts(quoted, limit)
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

// Problemas com a pesquisa atual:
// - Não consegue lidar com motorizações (1.0, 1.6), problema do SQLite
func serviceSearchCars(search string, limit int) (carSearchResponse, error) {
	fmt.Println("serviceSearchCars")
	fmt.Println("search:", search)
	fmt.Println("limit:", limit)
	if limit <= 0 {
		limit = 5
	} else if limit > 100 {
		limit = 100
	}

	search = strings.TrimSpace(search)
	r := strings.NewReplacer(
		" de ", " ",
		" da ", " ",
		" do ", " ")
	search = r.Replace(search)
	if search == "" {
		return carSearchResponse{}, errors.New("invalid query")
	}
	quoted := global.QuoteTokens([]string{search})[0]
	fmt.Println("service", quoted)
	rows, err := searchCars(quoted, limit)
	if err != nil {
		fmt.Println("service", err)
		return carSearchResponse{}, err
	}
	res, err := toCarSearchResponse(rows)
	if err != nil {
		fmt.Println("service", err)
		return carSearchResponse{}, err
	}
	return res, nil
}
