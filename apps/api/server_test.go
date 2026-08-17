package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	// 1. Crear una petición HTTP falsa simulando un método GET a /health
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("No se pudo crear la petición: %v", err)
	}

	// 2. Crear un ResponseRecorder para capturar la respuesta del servidor en memoria
	rr := httptest.NewRecorder()

	// 3. Ejecutar el handler directamente pasando el Recorder y la Request
	healthHandler(rr, req)

	// 4. Validar que el código de estado HTTP sea 200 OK
	if rr.Code != http.StatusOK {
		t.Errorf("Código de estado incorrecto: se esperaba %d, se obtuvo %d", http.StatusOK, rr.Code)
	}

	// 5. Validar que el Header Content-Type sea application/json
	expectedContentType := "application/json"
	if actual := rr.Header().Get("Content-Type"); actual != expectedContentType {
		t.Errorf("Content-Type incorrecto: se esperaba %q, se obtuvo %q", expectedContentType, actual)
	}

	// 6. Validar y decodificar el cuerpo JSON de la respuesta
	var response healthResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("No se pudo decodificar el JSON de respuesta: %v", err)
	}

	// 7. Verificar que el valor del campo "status" sea "up"
	expectedStatus := "up"
	if response.Status != expectedStatus {
		t.Errorf("Contenido del JSON incorrecto: se esperaba status=%q, se obtuvo status=%q", expectedStatus, response.Status)
	}
}
