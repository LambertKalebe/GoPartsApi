package search

type productSearchRequest struct {
	Search string `query:"search" example:"KL582"`
	Limit  int    `query:"limit" example:"1"`
}

type productSearchResponse struct {
	Products []product
}
