package export

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

func imagesExportHandler(c *echo.Context) error {
	req := new(productImageDownloadRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Corpo da requisição inválido")
	}
	urls, err := serviceGetProductImagesUrl(req.ProductId)
	if err != nil {
		c.JSON(500, err)
		return nil
	}

	//Não testado
	for i, url := range urls.ImageUrl {
		sep := "#"
		fileName := req.FileName + sep + strconv.Itoa(i) + ".jpg"

		return c.Attachment(url, fileName)
	}

	return nil
}
