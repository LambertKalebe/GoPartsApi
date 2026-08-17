package download

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
// @Success 200 {object} string
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

// @Summary Applications
// @Description Baixa as aplicações informadas no formato de csv
// @Tags Download
// @Produce json
// @Param id body appDownloadRequest true "Ids dos veiculos"
// @Router /export/apps [post]
// @Success 200 {object} string
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func appExportHandler(c *echo.Context) error {
	req := new(appDownloadRequest)

	if err := c.Bind(req); err != nil {
		return c.String(
			http.StatusBadRequest,
			"Corpo da requisição inválido",
		)
	}

	fmt.Printf("Handler Search: %#v\n", req.CarId)
	fmt.Printf("Handler Search Product: %#v\n", req.ProductId)
	fmt.Println("-----------------------------------")

	resp, err := serviceAppDownload(req.CarId, req.ProductId)
	if err != nil {
		return c.String(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	c.Response().Header().Set(
		"Content-Disposition",
		`attachment; filename="compat.csv"`,
	)

	fmt.Printf("Bytes enviados: % X\n", resp)

	return c.Blob(
		http.StatusOK,
		"text/csv;charset=utf-8",
		resp,
	)
}
