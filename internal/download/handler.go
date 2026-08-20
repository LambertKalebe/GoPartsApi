package download

import (
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// @Summary Images
// @Description Baixa todas as imagens de um produto
// @Tags Download
// @Produce image/png
// @Produce image/jpeg
// @Param fileName query string true "Nome base dos arquivos"
// @Param id query int true "ID do produto"
// @Param index query int true "Indice da imagem"
// @Router /go-api/export/images [get]
// @Success 200 {object} string
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func imagesExportHandler(c *echo.Context) error {
	req := new(productImageDownloadRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Corpo da requisição inválido")
	}

	urls, err := serviceGetProductImageUrl(req.ProductId)
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

	resp, err := http.Get(imageURL)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Erro ao baixar imagem")
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
// @Produce text/csv
// @Param id body appDownloadRequest true "Ids dos veiculos"
// @Router /go-api/export/apps [post]
// @Success 200 {file} file "Arquivo CSV"
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func appExportHandler(c *echo.Context) error {
	req := new(appDownloadRequest)

	if err := c.Bind(req); err != nil {
		return c.String(
			http.StatusBadRequest,
			"Corpo da requisição inválido",
		)
	}

	if req.CarId != nil && req.ProductId != 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest, "apenas um dos dois parâmetros pode ser selecionado: veiculo ou produto")
	}

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

	return c.Blob(
		http.StatusOK,
		"text/csv;charset=utf-8",
		resp,
	)
}
