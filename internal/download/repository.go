package download

import (
	"database/sql"
	"errors"
	"fmt"
	"g0/internal/database"
	"strings"
)

const getProductImagesUrlQuery = `
SELECT
    url
FROM product_images
WHERE product_id = ?
`
const getAppQueryByProdutID = `
SELECT     
    m.name AS montadora,     
    v.model AS veiculo,     
    v.version AS modelo,     
    e.code AS motor,     
    TRIM(         
        COALESCE(printf('%.1f', e.cc_normalized), '') ||         
        CASE             
            WHEN e.cc_normalized IS NOT NULL AND e.valves IS NOT NULL             
            THEN ' '             
            ELSE '0.0 0V'         
        END ||         
        COALESCE(e.valves, '')     ) AS config_motor,     
    t.name AS transmissao,     
    e.fuel AS combustivel,     
    v.year AS ano_inicio,     
    v.year AS ano_fim 
FROM vehicles v
JOIN product_applications pa ON pa.vehicle_id = v.id  -- assuming this is how product_applications links
JOIN makes m ON m.id = v.make_id 
LEFT JOIN engine_configs e ON e.id = v.engine_id 
LEFT JOIN transmissions t ON t.id = v.transmission_id 
WHERE pa.product_id = ?   
ORDER BY v.make_id;
`
const getAppQueryByVehicleIDs = `
SELECT
    m.name AS montadora,
    v.model AS veiculo,
    v.version AS modelo,
    COALESCE(e.code, '') AS motor,
    TRIM(
        COALESCE(printf('%%.1f', e.cc_normalized), '') ||
        CASE
            WHEN e.cc_normalized IS NOT NULL AND e.valves IS NOT NULL
            THEN ' '
            ELSE '0.0 0V'
        END ||
        COALESCE(e.valves, '')
    ) AS config_motor,
    COALESCE(t.name, '') AS transmissao,
    COALESCE(e.fuel, '') AS combustivel,
    COALESCE(v.year, 0) AS ano_inicio,
    COALESCE(v.year, 0) AS ano_fim
FROM vehicles v
JOIN makes m
    ON m.id = v.make_id
LEFT JOIN engine_configs e
    ON e.id = v.engine_id
LEFT JOIN transmissions t
    ON t.id = v.transmission_id
WHERE v.id IN (%s)
ORDER BY v.make_id;
`

func getProductImagesUrl(productId int) (*sql.Rows, error) {
	rows, err := database.DB.Query(getProductImagesUrlQuery, productId)
	if err != nil {
		return nil, err
	}

	return rows, nil

}

func getApps(vehicleIDs []int, productId int) (*sql.Rows, error) {
	if len(vehicleIDs) != 0 && productId != 0 {
		return nil, errors.New("apenas um dos dois parâmetros pode ser selecionado: veiculo ou produto")
	}

	if len(vehicleIDs) == 0 && productId == 0 {
		return nil, errors.New("nenhum veiculo ou produto foi selecionado")
	}
	var rows *sql.Rows
	var err error

	if productId != 0 {
		rows, err := database.DB.Query(getAppQueryByProdutID, productId)
		if err != nil {
			return nil, err
		}
		return rows, err
	}

	if len(vehicleIDs) != 0 {
		placeholders := make([]string, len(vehicleIDs))
		args := make([]any, len(vehicleIDs))

		for i, id := range vehicleIDs {
			placeholders[i] = "?"
			args[i] = id
		}

		query := fmt.Sprintf(
			getAppQueryByVehicleIDs,
			strings.Join(placeholders, ","),
		)

		rows, err = database.DB.Query(query, args...)
		if err != nil {
			return nil, err
		}
	}

	return rows, err

}
