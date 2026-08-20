package appbuilder

import "g0/internal/global"

type car struct {
	ID          int            `json:"ID" example:"20976"`
	Make        string         `json:"make" example:"VOLKSWAGEN"`
	Model       string         `json:"model" example:"Fox"`
	Version     string         `json:"version" example:"Plus 1.6 2009"`
	ConfigMotor string         `json:"configMotor" example:"1.6 8V"`
	Year        int            `json:"year" example:"2009"`
	FilterData  global.JsonMap `json:"filterData"`
}
