package vehicles

import "g0/internal/common"

type vehicle struct {
	ID               int    `json:"id"`
	Make             string `json:"make"`
	Model            string `json:"model"`
	Version          string `json:"version"`
	Year             int    `json:"year"`
	EngineCode       string `json:"engine"`
	ConfigMotor      string `json:"configMotor"`
	TransmissionType string `json:"transmission"`
	PartsCount       int    `json:"partsCount"`
}

type vehicleDetails struct {
	ID           int            `json:"id"`
	Make         string         `json:"make"`
	Model        string         `json:"model"`
	Version      string         `json:"version"`
	Year         int            `json:"year"`
	Engine       engine         `json:"engine"`
	Transmission transmission   `json:"transmission"`
	SourceUrl    string         `json:"source_url"`
	InfoData     common.JsonMap `json:"info"`
	TechData     common.JsonMap `json:"tech_data"`
	Parts        []parts        `json:"parts"`
}

type parts struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Make string `json:"make"`
	Name string `json:"name"`
}

type engine struct {
	Code       string `json:"code"`
	Valves     string `json:"valves"`
	CC         int    `json:"cc"`
	Fuel       string `json:"fuel"`
	Aspiration string `json:"aspiration"`
}

type transmission struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
