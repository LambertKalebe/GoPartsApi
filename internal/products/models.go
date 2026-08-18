package products

import "g0/internal/global"

type product struct {
	ID       int64  `json:"id" example:"1"`
	Code     string `json:"code" example:"FAP-2829"`
	Make     string `json:"make" example:"WEGA"`
	Name     string `json:"name" example:"FILTRO DE AR"`
	Public   int    `json:"public" example:"1"`
	ImageUrl string `json:"image" example:"https://www.wega.com.br/wp-content/uploads/2020/09/filtro-de-ar-wega-1.jpg"`
	AppCount int    `json:"appCount" example:"1"`
}

type productById struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Make     string `json:"make"`
	ImageUrl string `json:"image" example:"https://www.wega.com.br/wp-content/uploads/2020/09/filtro-de-ar-wega-1.jpg"`
	AppCount int    `json:"appCount" example:"1"`
}

type productDetails struct {
	ID           int64          `json:"id"`
	Code         string         `json:"code"`
	Name         string         `json:"name"`
	Make         string         `json:"make"`
	TechData     global.JsonMap `json:"tech_data"`
	LogisticData global.JsonMap `json:"logistic_data"`
	FiscalData   global.JsonMap `json:"fiscal_data"`
	Similar      []similar      `json:"similar"`
	Applications []application  `json:"applications"`
	Images       []image        `json:"images"`
}

type image struct {
	Number int    `json:"number"`
	URL    string `json:"url"`
}

type similar struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Make string `json:"make"`
}

type application struct {
	ID           int          `json:"id"`
	Make         string       `json:"make"`
	Model        string       `json:"model"`
	Version      string       `json:"version"`
	Year         int          `json:"year"`
	Engine       engine       `json:"engine"`
	Transmission transmission `json:"transmission"`
}

type engine struct {
	Code       string `json:"code"`
	Valves     string `json:"valves"`
	CC         string `json:"cc"`
	Fuel       string `json:"fuel"`
	Aspiration string `json:"aspiration"`
}

type transmission struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
