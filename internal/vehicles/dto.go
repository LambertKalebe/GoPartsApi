package vehicles

type vehiclesRequest struct {
	Limit int `query:"limit" example:"100"`
	Page  int `query:"page" example:"1"`
}

type vehicleByIdRequest struct {
	Id int `query:"id" example:"1"`
}

type vehiclesResponse struct {
	Page    int       `json:"page"`
	Vehicle []vehicle `json:"vehicle"`
}

type vehicleByIdResponse struct {
	Vehicle []vehicle `json:"vehicle"`
}

type vehicleDetailsByIdResponse struct {
	Vehicle []vehicleDetails `json:"vehicle"`
}
