package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Captura de rutas no registradas (404)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprintln(w, "Hello from, HTTP")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "not_found"})
	})

	mux.HandleFunc("GET /health", healthHandler)
	return mux
}
