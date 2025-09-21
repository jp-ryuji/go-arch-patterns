package http

import (
	"encoding/json"
	"net/http"
)

// HealthCheckHandler handles health check requests
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"message": "Service is running",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// In a production system, you would use a proper logger
		println("Failed to encode response:", err.Error())
	}
}

// NotFoundHandler handles 404 requests
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"error":   "Not found",
		"message": "The requested resource was not found",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// In a production system, you would use a proper logger
		println("Failed to encode response:", err.Error())
	}
}
