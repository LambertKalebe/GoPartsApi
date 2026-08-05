package health

func toHealthResponse(healthStatus bool, message string) healthResponse {
	return healthResponse{
		Healthy: healthStatus,
		Message: message,
	}
}
