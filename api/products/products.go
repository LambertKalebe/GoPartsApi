package products

import (
	"database/sql"
	"g0/internal/database"
	"net/http"

	"github.com/labstack/echo/v5"
)

type RequestProducts struct {
	Limit      int64 `query:"limit" example:"100"`
	PublicOnly bool  `query:"publicOnly" example:"true"`
}

type Product struct {
	ID       int64  `json:"id" example:"1"`
	Code     string `json:"code" example:"FAP-2829"`
	CodeNorm string `json:"codeNorm" example:"FAP2829"`
	Make     string `json:"make" example:"WEGA"`
	Name     string `json:"name" example:"FILTRO DE AR"`
	ImageUrl string `json:"image" example:"https://www.wega.com.br/wp-content/uploads/2020/09/filtro-de-ar-wega-1.jpg"`
	Public   int    `json:"public" example:"1"`
	AppCount int    `json:"appCount" example:"1"`
}

type ResponseProducts struct {
	Products []Product `json:"products"`
}

const GetProductsQuery = `
SELECT
    p.id,
    p.code,
    p.code_norm,
    m.name,
    p.name,
    (
        SELECT COUNT(*)
        FROM product_applications pa
        WHERE pa.product_id = p.id
    ) AS app_count,
    (
        SELECT product_images.url
        FROM product_images
        WHERE product_images.product_id = p.id
        LIMIT 1
    ) as image,
    p.public
FROM (
    SELECT *
    FROM products
    WHERE (? = 0 OR public = 1)
    ORDER BY id
    LIMIT ?
) p
LEFT JOIN makes m
    ON m.id = p.make_id;
`

func Route(g *echo.Group) {
	g.GET("", Products)
}

// Products
// @Summary Products
// @Description Realiza uma consulta de produtos (id 1 ao `limit`)
// @Tags Products
// @Produce json
// @Param limit query int false "Quantidade máxima de produtos"
// @Param publicOnly query bool false "Somente produtos públicos"
// @Router /products [get]
// @Success 200 {object} ResponseProducts
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func Products(c *echo.Context) error {
	req := new(RequestProducts)

	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Corpo da requisição inválido")
	}

	if req.Limit <= 0 {
		req.Limit = 100
	}

	filter := 0
	if req.PublicOnly {
		filter = 1
	}

	rows, err := database.DB.Query(GetProductsQuery, filter, req.Limit)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	resp := ResponseProducts{
		Products: []Product{},
	}

	for rows.Next() {
		var (
			p     Product
			image sql.NullString
		)

		err := rows.Scan(
			&p.ID,
			&p.Code,
			&p.CodeNorm,
			&p.Make,
			&p.Name,
			&p.AppCount,
			&image,
			&p.Public,
		)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		if image.Valid {
			p.ImageUrl = image.String
		}

		resp.Products = append(resp.Products, p)
	}

	if err := rows.Err(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
