package download

import (
	"database/sql"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

	defer rows.Close()

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

func removerAcentos(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	resultado, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return resultado
}

func (v *vehicle) limparAcentos() {
	v.Montadora = removerAcentos(v.Montadora)
	v.Veiculo = removerAcentos(v.Veiculo)
	v.Modelo = removerAcentos(v.Modelo)
	v.Motor = removerAcentos(v.Motor)
	v.ConfigMotor = removerAcentos(v.ConfigMotor)
	v.Transmissao = removerAcentos(v.Transmissao)
	v.Combustivel = removerAcentos(v.Combustivel)
}
