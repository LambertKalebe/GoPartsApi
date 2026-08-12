package export

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
