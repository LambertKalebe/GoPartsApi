package products

type productRequest struct {
	Qnt        int  `query:"limit" example:"100"`
	Page       int  `query:"page" example:"1"`
	PublicOnly bool `query:"publicOnly" example:"true"`
}

type productsResponse struct {
	Page     int       `json:"page"`
	Products []product `json:"productsHandler"`
}

type productByIdResponse struct {
	Product []productById `json:"product"`
}

type productDetailsByIdResponse struct {
	Product []productDetails `json:"productsHandler"`
}
