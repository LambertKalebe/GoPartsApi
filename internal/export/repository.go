package export

import (
	"database/sql"
	"g0/internal/database"
)

const getProductImagesUrlQuery = `
SELECT
    url
FROM product_images
WHERE product_id = ?
`

func getProductImagesUrl(productId int) (*sql.Rows, error) {
	rows, err := database.DB.Query(getProductImagesUrlQuery, productId)
	if err != nil {
		return nil, err
	}

	return rows, nil

}
