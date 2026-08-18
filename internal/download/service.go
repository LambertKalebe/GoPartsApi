package download

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/gocarina/gocsv"
)

func serviceGetProductImageUrl(productId int) (productImageQueryResponse, error) {
	req := productImageDownloadRequest{
		ProductId: productId,
	}
	urls, err := getProductImagesUrl(req.ProductId)
	if err != nil {
		return productImageQueryResponse{}, err
	}

	res, err := toProductImageDownloadResponse(urls)
	if err != nil {
		return productImageQueryResponse{}, err
	}

	return res, nil
}

func serviceAppDownload(vehicleIDs []int, productId int) ([]byte, error) {
	rows, err := getApps(vehicleIDs, productId)
	if err != nil {
		return nil, err
	}
	vehicles, err := formatSqlAppResponse(rows)
	// Configuração do CSV
	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		// Adiciona o BOM, por mais que não seja necessário, foi adicionado devido a outro setor da empresa que pediu
		_, err := out.Write([]byte{0xEF, 0xBB, 0xBF})
		if err != nil {
			return nil
		}
		// Adiciona o BOM novamente. Aparentemente os browsers removem o BOM quando o arquivo é baixado.
		_, err2 := out.Write([]byte{0xEF, 0xBB, 0xBF})
		if err2 != nil {
			return nil
		}

		writer := csv.NewWriter(out)
		writer.UseCRLF = true
		writer.Comma = ';'
		return gocsv.NewSafeCSVWriter(writer)
	})
	if err != nil {
		return nil, err
	}

	// CSV em UTF-8
	var utf8Buffer bytes.Buffer

	if err := gocsv.Marshal(&vehicles, &utf8Buffer); err != nil {
		return nil, fmt.Errorf("erro ao gerar CSV: %w", err)
	}

	return utf8Buffer.Bytes(), nil
}
