package productIdInfo

import (
	"database/sql"
	"encoding/json"
	"g0/internal/database"
	"net/http"

	"github.com/labstack/echo/v5"
)

func Route(g *echo.Group) {
	g.GET("/:id", productInfo)
}

type product struct {
	Id           int    `json:"id" example:"1"`
	Code         string `json:"code" example:"FAP-2829"`
	Name         string `json:"name" example:"FILTRO DE AR"`
	Make         string `json:"make" example:"WEGA"`
	TechData     any    `json:"tech_data"`
	LogisticData any    `json:"logistic_data"`
	FiscalData   any    `json:"fiscal_data"`
	SimilarJson  any    `json:"similar_json"`
	Application  any    `json:"application"`
	Images       any    `json:"images"`
}
type response struct {
	Products []product `json:"products"`
}

const getProductInfoQuery = `
SELECT
    p.id,
    p.code,
    p.name,
    m.name AS make,
    p.tech_data,
    p.logistic_data,
    p.fiscal_data,

    -- Similares
    (
        SELECT json_group_array(
            json_object(
                'id', pc.id,
                'code', pc.equivalent_code,
                'make', sm.name
            )
        )
        FROM product_crossrefs pc
        LEFT JOIN makes sm
            ON sm.id = pc.equivalent_make_id
        WHERE pc.product_id = p.id
    ) AS similar_json,

    -- Aplicações
    (
        SELECT json_group_array(
            json_object(
                'make', vm.name,
                'model', v.model,
                'version', v.version,
                'year', v.year,

                'engine', json_object(
                    'code', ec.code,
                    'valves', ec.valves,
                    'cc', ec.cc,
                    'fuel', ec.fuel,
                    'aspiration', ec.aspiration
                ),

                'transmission', json_object(
                    'code', t.code,
                    'name', t.name
                )
            )
        )
        FROM product_applications pa
        JOIN vehicles v
            ON v.id = pa.vehicle_id
        LEFT JOIN makes vm
            ON vm.id = v.make_id
        LEFT JOIN engine_configs ec
            ON ec.id = v.engine_id
        LEFT JOIN transmissions t
            ON t.id = v.transmission_id
        WHERE pa.product_id = p.id
    ) AS applications,

    -- Imagens
    (
        SELECT json_group_array(
            json_object(
                'id', pi.id,
                'number', pi.number,
                'url', pi.url
            )
        )
        FROM product_images pi
        WHERE pi.product_id = p.id
    ) AS images

FROM products p

LEFT JOIN makes m
    ON m.id = p.make_id

WHERE p.id = ?;
`

// ProductInfo
//
//	@Summary					ProductInfo
//	@Description				Retorna todos os dados do produto
//	@Tags						Products
//	@Security					CookieAuth
//	@securityDefinitions.apikey	CookieAuth
//	@in							cookie
//	@name						token
//	@Produce					json
//	@Success					200	{object}	response
//	@Failure					401	{object}	error
//	@Router						/products/{id} [get]
func productInfo(c *echo.Context) error {
	id := c.Param("id")

	rows, err := database.DB.Query(getProductInfoQuery, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	resp := response{
		Products: []product{},
	}

	for rows.Next() {
		var (
			p           product
			image       sql.NullString
			fiscal      string
			logistic    string
			similar     string
			technical   string
			application string
		)
		err := rows.Scan(
			&p.Id,
			&p.Code,
			&p.Name,
			&p.Make,
			&technical,
			&logistic,
			&fiscal,
			&similar,
			&application,
			&image,
		)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		p.FiscalData = json.RawMessage(fiscal)
		p.LogisticData = json.RawMessage(logistic)
		p.SimilarJson = json.RawMessage(similar)
		p.TechData = json.RawMessage(technical)
		p.Application = json.RawMessage(application)

		if image.Valid {
			p.Images = json.RawMessage(image.String)
		}

		resp.Products = append(resp.Products, p)
	}

	if err := rows.Err(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
