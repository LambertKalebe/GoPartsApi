package export

type productImageDownloadRequest struct {
	FileName  string `query:"fileName" example:"000001"`
	ProductId int    `query:"id" example:"1"`
	Index     int    `query:"index" example:"1"`
}

type productImageQueryResponse struct {
	ImageUrl []string
}

type productImageDownloadResponse struct {
	Image string
}
