package health

import (
	"g0/internal/database"
)

func serviceHealth() healthResponse {
	if database.Connect() != nil {
		return toHealthResponse(false, "OK")
	}
	return toHealthResponse(true, "OK")
}
