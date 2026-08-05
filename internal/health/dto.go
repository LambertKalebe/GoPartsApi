package health

type healthResponse struct {
	Healthy bool
	Message string `json:"message" example:"On"`
}
