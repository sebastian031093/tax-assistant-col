package server

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status   string `json:"status"`
	Services string `json:"service"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := healthResponse{
		Status:   "ok",
		Services: "tax-assistant-api",
	}
	json.NewEncoder(w).Encode(response)
}
