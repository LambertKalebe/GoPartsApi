package download

import (
	"database/sql"
)

func toProductImageDownloadResponse(rows *sql.Rows) (productImageQueryResponse, error) {
	resp := productImageQueryResponse{}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var i image

		_ = rows.Scan(
			&i.URL,
		)

		resp.ImageUrl = append(resp.ImageUrl, i.URL)
	}
	return resp, nil
}

func formatSqlAppResponse(rows *sql.Rows) ([]vehicle, error) {
	res := make([]vehicle, 0)

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var v vehicle

		err := rows.Scan(
			&v.Montadora,
			&v.Veiculo,
			&v.Modelo,
			&v.Motor,
			&v.ConfigMotor,
			&v.Transmissao,
			&v.Combustivel,
			&v.AnoInicio,
			&v.AnoFim,
		)
		if err != nil {
			return nil, err
		}

		v.Obs = ""
		res = append(res, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
