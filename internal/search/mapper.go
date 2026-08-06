package search

import "database/sql"

func toProductSearchResponse(rows *sql.Rows) (productSearchResponse, error) {
	resp := productSearchResponse{}
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
		)

		if image.Valid {
			p.ImageUrl = image.String
		}

		resp.Products = append(resp.Products, p)
	}
	return resp, nil
}
