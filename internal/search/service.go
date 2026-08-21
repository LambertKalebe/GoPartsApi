package search

import (
	"errors"
	"g0/internal/common"
	"strings"
)

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
	search = r.Replace(strings.ToLower(search))
	if search == "" {
		return productSearchResponse{}, errors.New("invalid query")
	}
	quoted := common.QuoteTokens([]string{search})[0]
	rows, err := searchProducts(quoted, search, limit)
	if err != nil {
		return productSearchResponse{}, err
	}
	res, err := toProductSearchResponse(rows)
	if err != nil {
		return productSearchResponse{}, err
	}
	return res, nil
}

func serviceSearchCars(search string, limit int) (carSearchResponse, error) {
	if limit <= 0 {
		limit = 5
	} else if limit > 100 {
		limit = 100
	}

	search = strings.TrimSpace(search)
	r := strings.NewReplacer(
		" de ", " ",
		" da ", " ",
		" do ", " ",
		"vw", "volkswagen",
		"gm", "chevrolet")
	search = r.Replace(strings.ToLower(search))
	if search == "" {
		return carSearchResponse{}, errors.New("invalid query")
	}
	quoted := common.QuoteTokens([]string{search})[0]
	rows, err := searchCars(quoted, limit)
	if err != nil {
		return carSearchResponse{}, err
	}
	res, err := toCarSearchResponse(rows)
	if err != nil {
		return carSearchResponse{}, err
	}
	return res, nil
}
