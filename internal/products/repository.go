package products

import (
	"database/sql"
	"g0/internal/database"
)

const getProductsQuery = `
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
        LIMIT 1
    ) AS image,
    p.public
FROM products p
LEFT JOIN makes m
    ON m.id = p.make_id
WHERE
    (? = 0 OR p.public = 1)
ORDER BY p.id
LIMIT ?
OFFSET ?;
`

const getProductByIdQuery = `
SELECT
    p.id,
    p.code,
    p.name,
    m.name AS make,
    (
        SELECT COUNT(*)
        FROM product_applications pa
        WHERE pa.product_id = p.id
    ) AS app_count,
    (
        SELECT pi.url
        FROM product_images pi
        WHERE pi.product_id = p.id
        LIMIT 1
    ) AS image
FROM products p
LEFT JOIN makes m
    ON m.id = p.make_id
WHERE p.id = ?;
`

const getProductDetailByIdQuery = `
SELECT
    p.id,
    p.code,
    p.name,
    m.name AS make,
    p.tech_data,
    p.logistic_data,
    p.fiscal_data
FROM products p
LEFT JOIN makes m
    ON m.id = p.make_id
WHERE p.id = ?;
`

const getProductCrossrefsByProductIdQuery = `
SELECT
    pc.equivalent_product_id,
    pc.equivalent_code AS code,
    sm.name AS make
FROM product_crossrefs pc
LEFT JOIN makes sm
    ON sm.id = pc.equivalent_make_id
WHERE pc.product_id = ?;
`

const getProductsAppsByProductIdQuery = `
SELECT
    v.id,
    vm.name AS make,
    v.model,
    v.version,
    v.year,

    ec.code AS engine_code,
    ec.valves,
    printf('%.1f', ec.cc_normalized) AS cc,
    ec.fuel,
    ec.aspiration,

    t.code AS transmission_code,
    t.name AS transmission_name

FROM product_applications pa

JOIN vehicles v
    ON v.id = pa.vehicle_id

LEFT JOIN makes vm
    ON vm.id = v.make_id

LEFT JOIN engine_configs ec
    ON ec.id = v.engine_id

LEFT JOIN transmissions t
    ON t.id = v.transmission_id

WHERE pa.product_id = ?;
`

const getProductImagesByProductIdQuery = `
SELECT
    number,
    url
FROM product_images
WHERE product_id = ?;
`

func getProducts(qnt int, page int, public bool) (*sql.Rows, error) {

	if qnt <= 0 {
		qnt = 100
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * qnt

	filter := 0

	if public {
		filter = 1
	}

	rows, err := database.DB.Query(getProductsQuery, filter, qnt, offset)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func getProductById(id int) (*sql.Row, error) {
	row := database.DB.QueryRow(
		getProductByIdQuery,
		id,
	)
	return row, nil
}

//----------------------Details--------------------------

func getProductDetailsById(id int) (*sql.Row, error) {
	row := database.DB.QueryRow(getProductDetailByIdQuery, id)
	return row, nil
}

func getProductCrossrefsByProductId(id int) (*sql.Rows, error) {
	rows, err := database.DB.Query(getProductCrossrefsByProductIdQuery, id)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func getProductsAppsByProductId(id int) (*sql.Rows, error) {
	rows, err := database.DB.Query(getProductsAppsByProductIdQuery, id)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func getProductImagesByProductId(id int) (*sql.Rows, error) {
	rows, err := database.DB.Query(getProductImagesByProductIdQuery, id)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
