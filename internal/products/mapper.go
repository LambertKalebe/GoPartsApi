package products

import (
	"database/sql"
)

func toProductsResponse(rows *sql.Rows, page int) (productsResponse, error) {

	resp := productsResponse{}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var p product
		var image sql.NullString

		_ = rows.Scan(
			&p.ID,
			&p.Code,
			&p.Make,
			&p.Name,
			&p.AppCount,
			&image,
			&p.Public,
		)

		if image.Valid {
			p.ImageUrl = image.String
		}

		resp.Products = append(resp.Products, p)
	}
	resp.Page = page
	return resp, nil
}

func toProductByIdResponse(row *sql.Row) (productByIdResponse, error) {
	resp := productByIdResponse{}

	var p productById
	var image sql.NullString

	err := row.Scan(
		&p.ID,
		&p.Code,
		&p.Make,
		&p.Name,
		&p.AppCount,
		&image,
	)
	if err != nil {
		return resp, err
	}

	if image.Valid {
		p.ImageUrl = image.String
	}

	resp.Product = append(resp.Product, p)
	return resp, nil
}

func toProductDetailsResponse(productDetailsRow *sql.Row, productImagesRows *sql.Rows, productCrossrefRows *sql.Rows, productAppRows *sql.Rows) (productDetailsByIdResponse, error) {

	var p productDetails
	var img image
	var s similar
	var app application

	resp := productDetailsByIdResponse{}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(productImagesRows)

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(productCrossrefRows)

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(productAppRows)

	for productImagesRows.Next() {
		err := productImagesRows.Scan(
			&img.Number,
			&img.URL,
		)
		if err != nil {
			return resp, err
		}

		p.Images = append(p.Images, img)
	}

	for productCrossrefRows.Next() {
		err := productCrossrefRows.Scan(
			&s.ID,
			&s.Code,
			&s.Make,
		)

		if err != nil {
			return resp, err
		}

		p.Similar = append(p.Similar, s)
	}

	for productAppRows.Next() {
		err := productAppRows.Scan(
			&app.ID,
			&app.Make,
			&app.Model,
			&app.Version,
			&app.Year,
			&app.Engine.Code,
			&app.Engine.Valves,
			&app.Engine.CC,
			&app.Engine.Fuel,
			&app.Engine.Aspiration,
			&app.Transmission.Code,
			&app.Transmission.Name,
		)

		if err != nil {
			return resp, err
		}

		p.Applications = append(p.Applications, app)
	}
	err := productDetailsRow.Scan(
		&p.ID,
		&p.Code,
		&p.Name,
		&p.Make,
		&p.TechData,
		&p.LogisticData,
		&p.FiscalData,
	)

	if err != nil {
		return resp, err
	}

	resp.Product = append(resp.Product, p)

	return resp, nil

}
