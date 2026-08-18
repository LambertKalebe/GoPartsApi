package search

import (
	"database/sql"
	"g0/internal/database"
)

const productSearchQuery = `
SELECT
    p.id,
    p.code,
    m.name AS make,
    p.name,
    (
        SELECT COUNT(*)
        FROM product_applications pa
        WHERE pa.product_id = p.id
    ) AS app_count,
    (
        SELECT pi.url
        FROM product_images pi
        WHERE pi.product_id = p.id
        ORDER BY pi.number
        LIMIT 1
    ) AS image
FROM fts5_products f
JOIN products p
    ON p.id = f.rowid
LEFT JOIN makes m
    ON m.id = p.make_id
WHERE
    fts5_products MATCH ?
    AND p.public = 1
ORDER BY
    CASE
        WHEN LOWER(p.code_norm) = LOWER(?)
          OR LOWER(p.code) = LOWER(?)
        THEN 0
        ELSE 1
    END,
    app_count DESC,
    p.name
LIMIT ?;
`

const carSearchQuery = `
SELECT
    v.id,
    m.name AS make,
    v.model,
    v.version,
    v.year,
    e.code AS engineCode,
    t.code AS transmissionCode
FROM fts5_vehicles f
JOIN vehicles v
    ON v.id = f.rowid
LEFT JOIN makes m
    ON m.id = v.make_id
LEFT JOIN engine_configs e
    ON e.id = v.engine_id
LEFT JOIN transmissions t
    ON t.id = v.transmission_id
WHERE fts5_vehicles MATCH ?
LIMIT ?;
`

func searchProducts(search string, searchBase string, limit int) (*sql.Rows, error) {

	rows, err := database.DB.Query(
		productSearchQuery,
		search,
		searchBase,
		searchBase,
		limit,
	)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func searchCars(search string, limit int) (*sql.Rows, error) {

	rows, err := database.DB.Query(
		carSearchQuery,
		search,
		limit,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
