package search

type productSearchRequest struct {
	Query string
}

type productSearchResponse struct {
	Products []product
}
