package appbuilder

type car struct {
	ID          int    `json:"id"`
	Make        string `json:"make"`
	Model       string `json:"model"`
	Version     string `json:"version"`
	ConfigMotor string `json:"configMotor"`
	Year        int    `json:"year"`
}
