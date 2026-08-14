package download

type productImageDownloadRequest struct {
	FileName  string `query:"fileName" example:"000001"`
	ProductId int    `query:"id" example:"1"`
	Index     int    `query:"index" example:"1"`
}

type productImageQueryResponse struct {
	ImageUrl []string
}

type appDownloadRequest struct {
	ProductId int   `json:"productId" example:"1"`
	CarId     []int `json:"vehicleId" example:"1"`
}

type vehicle struct {
	Montadora   string `json:"montadora" csv:"Montadora"`
	Veiculo     string `json:"veiculo" csv:"Veiculo"`
	Modelo      string `json:"modelo" csv:"Modelo"`
	Motor       string `json:"motor" csv:"Motor"`
	ConfigMotor string `json:"configMotor" csv:"ConfigMotor"`
	Transmissao string `json:"transmissao" csv:"Transmissao"`
	Combustivel string `json:"combustivel" csv:"Combustivel"`
	AnoInicio   int    `json:"anoInicio" csv:"AnoInicio"`
	AnoFim      int    `json:"anoFim" csv:"AnoFim"`
	Obs         string `json:"obs" csv:"Obs"`
}
