package vehicles

import (
	"database/sql"
	"g0/internal/database"
)

const getVehiclesQuery = `
SELECT
    v.id,
    m.name AS make,
    v.model,
    v.version,
    e.code AS engineCode,
    t.code AS transmissionCode,
    v.year,
    (
        SELECT COUNT(*)
        FROM product_applications pa
        WHERE pa.vehicle_id = v.id
    ) AS partsCount
FROM vehicles v
LEFT JOIN makes m
    ON m.id = v.make_id
LEFT JOIN engine_configs e
    ON e.id = v.engine_id
LEFT JOIN transmissions t
    ON t.id = v.transmission_id
ORDER BY v.id
LIMIT ?
OFFSET ?;
`

const getVehicleByIdQuery = `
SELECT
    v.id,
    m.name AS make,
    v.model,
    v.version,
    e.code AS engineCode,
    t.code AS transmissionCode,
    v.year,
    (
        SELECT COUNT(*)
        FROM product_applications pa
        WHERE pa.vehicle_id = v.id
    ) AS partsCount
FROM vehicles v
LEFT JOIN makes m
    ON m.id = v.make_id
LEFT JOIN engine_configs e
    ON e.id = v.engine_id
LEFT JOIN transmissions t
    ON t.id = v.transmission_id
WHERE v.id = ?;
`
const getVehicleDetailByIdQuery = `
SELECT
    v.id,
    m.name AS make,
    v.model,
    v.version,
    v.year,
    e.code AS engine_code,
    e.valves AS engine_valves,
    COALESCE(NULLIF(e.cc, ''), 0) AS engine_cc,
    e.fuel AS engine_fuel,
    e.aspiration AS engine_aspiration,
    t.code AS transmission_code,
    t.name AS transmission_name,
    vp.info_json,
    vp.tech_json,
    v.source_url
	
FROM vehicles v
LEFT JOIN makes m
    ON m.id = v.make_id
LEFT JOIN vehicle_profiles vp
    ON vp.vehicle_id = v.id
LEFT JOIN engine_configs e
    ON e.id = v.engine_id
LEFT JOIN transmissions t
    ON t.id = v.transmission_id
WHERE v.id = ?;
`

const getVehicleAppsByIdQuery = `
SELECT
    p.id,
    p.name,
    m.name AS make,
    p.code
FROM product_applications pa
JOIN products p
    ON p.id = pa.product_id
LEFT JOIN makes m
    ON m.id = p.make_id
WHERE pa.vehicle_id = ?
ORDER BY p.id;
`

func getVehicles(qnt int, page int) (*sql.Rows, error) {

	if qnt <= 0 {
		qnt = 100
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * qnt

	rows, err := database.DB.Query(getVehiclesQuery, qnt, offset)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func getVehicleById(id int) (*sql.Row, error) {
	row := database.DB.QueryRow(
		getVehicleByIdQuery,
		id,
	)
	return row, nil
}

//----------------------Details--------------------------

func getVehicleDetailsById(id int) (*sql.Row, error) {
	row := database.DB.QueryRow(getVehicleDetailByIdQuery, id)
	return row, nil
}

func getVehicleAppsByProductId(id int) (*sql.Rows, error) {
	rows, err := database.DB.Query(getVehicleAppsByIdQuery, id)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
