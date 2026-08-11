package export

import (
	"database/sql"
	"g0/internal/database"
)

const getProductImagesUrlQuery = `
SELECT
    id,
    url
FROM product_images
WHERE product_id = ?
`

func getProductImagesUrl(productId int) (*sql.Rows, error) {
	rows, err := database.DB.Query(getProductImagesUrlQuery, id)
	if err != nil {
		return nil, err
	}

	return rows, nil

}
