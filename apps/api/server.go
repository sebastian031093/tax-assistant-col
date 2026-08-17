package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// JSON response structure
type healthResponse struct {
	Status   string `json:"status"`
	Services string `json:"service"`
}

// 1. Handler function for the health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header so clients know to parse it as JSON
	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK) // 200 OK

	response := healthResponse{
		Status:   "ok",
		Services: "tax-assistant-api",
	}

	// Encode the struct directly into the response writer
	json.NewEncoder(w).Encode(response)
}

// NUEVA FUNCIÓN: Fábrica del Enrutador (Router)
// Esta función centraliza todas las rutas de tu aplicación.
func newRouter() *http.ServeMux {
	// Create a dedicated request multiplexer (router)
	mux := http.NewServeMux()
	// Global / catch-all route (plain text)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" {
			fmt.Fprintln(w, "Hello from, HTTP")
		}
		// Si el path NO es "/", significa que es una ruta no registrada (ej: /unknown)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound) // 404

		// Respondemos con un JSON personalizado de error
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "not_found",
			"message": "La ruta solicitada no existe.",
		})
	})
	// GET routing: Using "GET " prefix strictly limits this path to GET requests
	mux.HandleFunc("GET /health", healthHandler)

	return mux
}

// 2. Server Configuration factory
func newServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:         ":8080",
		Handler:      mux,              // Use our configured router instead of nil
		ReadTimeout:  5 * time.Second,  // Max time to read request
		WriteTimeout: 10 * time.Second, // Max time to write response
		IdleTimeout:  15 * time.Second, // Max time to keep connection alive
	}
}

// 3. Startup and Routing Initialization
func main() {

	fmt.Println("Hello world from Go, Sebastian")
	//start: app rounting
	mux := newRouter()
	// Initialize the custom server configuration
	server := newServer(mux)

	log.Printf("starting server on %s", server.Addr)

	// Start the server using the custom configuration
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
