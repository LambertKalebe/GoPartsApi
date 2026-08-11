package export

func serviceGetProductImagesUrl(productId int) (productImageDownloadResponse, error) {
	req := productImageDownloadRequest{
		ProductId: productId,
	}
	url, err := getProductImagesUrl(req.ProductId)
	if err != nil {
		return productImageDownloadResponse{}, err
	}

	res, err := toProductImageDownloadResponse(url)
	if err != nil {
		return productImageDownloadResponse{}, err
	}

	return res, nil
}
