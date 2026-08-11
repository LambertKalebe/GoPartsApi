package export

type productImageDownloadRequest struct {
	FileName  string
	ProductId int
}

type productImageDownloadResponse struct {
	ImageUrl []string
}
