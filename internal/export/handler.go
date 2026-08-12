package export

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// @Summary Images
// @Description Baixa todas as imagens de um produto
// @Tags Download
// @Produce json
// @Param fileName query string true "Nome base dos arquivos"
// @Param id query int true "ID do produto"
// @Param index query int true "Indice da imagem"
// @Router /export/images [get]
// @Success 200 {object} productImageDownloadResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func imagesExportHandler(c *echo.Context) error {
	req := new(productImageDownloadRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Corpo da requisição inválido")
	}
	fmt.Println("Handler REQ", req)
	urls, err := serviceGetProductImageUrl(req.ProductId)
	fmt.Println("Handler URLS", urls)
	if err != nil {
		err := c.JSON(500, err)
		if err != nil {
			return err
		}
		return nil
	}

	if req.Index >= len(urls.ImageUrl) {
		err := c.JSON(404, "Imagem não encontrada")
		if err != nil {
			return err

		}
		return nil
	}
	imageURL := urls.ImageUrl[req.Index]

	fmt.Println("Handler URL", imageURL)
	fmt.Println("Handler TOTAL IMAGES", len(urls.ImageUrl))
	resp, err := http.Get(imageURL)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao baixar imagem")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	sep := "#"
	fileName := req.FileName + sep + strconv.Itoa(req.Index) + ".jpg"

	c.Response().Header().Set(
		"Content-Disposition",
		`attachment; filename="`+fileName+`"`,
	)

	c.Response().Header().Set("X-Images-Count",
		strconv.Itoa(len(urls.ImageUrl)))

	return c.Stream(
		http.StatusOK,
		"image/jpeg",
		resp.Body,
	)

}
