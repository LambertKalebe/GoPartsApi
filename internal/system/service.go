package system

import (
	"g0/internal/database"
)

func serviceHealth() healthResponse {
	if database.Connect() != nil {
		return toHealthResponse(false, "Database connection failed")
	}
	return toHealthResponse(true, "OK")
}
