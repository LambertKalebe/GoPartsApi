package search

import (
	"database/sql"
	"fmt"
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
    app_count DESC,
    p.name
LIMIT ?;
`

func searchProducts(search string, limit int) (*sql.Rows, error) {

	fmt.Println("search:", search)
	fmt.Println("limit:", limit)

	rows, err := database.DB.Query(
		productSearchQuery,
		search,
		limit,
	)
	if err != nil {
		fmt.Println("Repository:", err)
		return nil, err
	}

	return rows, nil
}
