package products

import (
	"database/sql"
	"errors"
)

func serviceProducts(qnt int, page int, public bool) (productsResponse, error) {
	rows, err := getProducts(qnt, page, public)
	if err != nil {
		return productsResponse{}, err
	}
	res, err := toProductsResponse(rows, page)
	if err != nil {
		return productsResponse{}, err
	}
	if len(res.Products) == 0 {
		return productsResponse{}, errors.New("no products found")
	}

	return res, nil

}

func serviceProductById(id int) (productByIdResponse, error) {
	rows, err := getProductById(id)
	if err != nil {
		return productByIdResponse{}, err
	}
	res, err := toProductByIdResponse(rows)
	if err != nil {
		return productByIdResponse{}, err
	}
	return res, nil
}

func serviceProductDetailsById(id int) (productDetailsByIdResponse, error) {

	if id <= 0 {
		return productDetailsByIdResponse{}, errors.New("invalid query")
	}

	data, err := getProductDetailsById(id)
	if errors.Is(err, sql.ErrNoRows) {
		return productDetailsByIdResponse{}, errors.New("invalid query")
	}

	images, err := getProductImagesByProductId(id)
	if err != nil {
		return productDetailsByIdResponse{}, err
	}

	crossrefs, err := getProductCrossrefsByProductId(id)
	if err != nil {
		return productDetailsByIdResponse{}, err
	}

	apps, err := getProductsAppsByProductId(id)
	if err != nil {
		return productDetailsByIdResponse{}, err
	}

	res, err := toProductDetailsResponse(data, images, crossrefs, apps)
	if err != nil {
		return productDetailsByIdResponse{}, err
	}
	return res, nil
}
