package export

import "database/sql"

func toProductImageDownloadResponse(rows *sql.Rows) (productImageDownloadResponse, error) {
	resp := productImageDownloadResponse{}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var i image

		_ = rows.Scan(
			&i.ID,
			&i.URL,
		)

		resp.ImageUrl = append(resp.ImageUrl, i.URL)
	}
	return resp, nil
}
