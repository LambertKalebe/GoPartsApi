package search

type product struct {
	ID       int64  `json:"id" example:"1"`
	Code     string `json:"code" example:"FAP-2829"`
	Make     string `json:"make" example:"WEGA"`
	Name     string `json:"name" example:"FILTRO DE AR"`
	ImageUrl string `json:"image" example:"https://www.wega.com.br/wp-content/uploads/2020/09/filtro-de-ar-wega-1.jpg"`
	AppCount int    `json:"appCount" example:"1"`
}
