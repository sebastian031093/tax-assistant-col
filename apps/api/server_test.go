package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	// 1. Inicializamos el router real de la aplicación para probar el enrutamiento completo
	router := newRouter()

	// 2. Definición de la tabla de pruebas (Table-Driven Test)
	tests := []struct {
		name         string         // Nombre descriptivo del caso de uso
		method       string         // Método HTTP a simular
		path         string         // Ruta a la que apuntará la petición
		wantStatus   int            // Código de estado HTTP esperado
		validateBody bool           // Bandera para saber si este caso debe validar el JSON
		expectedBody healthResponse // Estructura esperada en caso de validar el JSON
	}{
		{
			name:         "GET /health exitoso",
			method:       "GET",
			path:         "/health",
			wantStatus:   http.StatusOK, // 200
			validateBody: true,
			expectedBody: healthResponse{
				Status:   "ok",
				Services: "tax-assistant-api",
			},
		},
		{
			name:         "POST /health no permitido",
			method:       "POST",
			path:         "/health",
			wantStatus:   http.StatusMethodNotAllowed, // 405 debido al prefijo "GET " en el mux
			validateBody: false,
		},
		{
			name:         "GET /unknown ruta inexistente",
			method:       "GET",
			path:         "/unknown",
			wantStatus:   http.StatusNotFound, // 404
			validateBody: false,
		},
	}

	// 3. Iteración sobre cada caso de la tabla
	for _, tc := range tests {
		// t.Run nos permite ejecutar sub-pruebas aisladas y reportar fallos específicos por nombre
		t.Run(tc.name, func(t *testing.T) {

			// Crear la petición HTTP falsa basada en los datos del caso actual
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatalf("No se pudo crear la petición: %v", err)
			}

			// Crear el ResponseRecorder para capturar la respuesta en memoria
			rr := httptest.NewRecorder()

			// IMPORTANTE: En lugar de llamar al handler directamente, le pasamos la petición
			// al router mediante ServeHTTP. Esto simula el ciclo de vida real del servidor.
			router.ServeHTTP(rr, req)

			// Validar el código de estado HTTP obtenido contra el esperado
			if rr.Code != tc.wantStatus {
				t.Errorf("Código de estado incorrecto: se esperaba %d, se obtuvo %d", tc.wantStatus, rr.Code)
			}

			// Validar el cuerpo y las cabeceras JSON solo si el caso de prueba lo requiere
			if tc.validateBody {
				// Validar que el Header Content-Type sea application/json
				expectedContentType := "application/json"
				if actual := rr.Header().Get("Content-Type"); actual != expectedContentType {
					t.Errorf("Content-Type incorrecto: se esperaba %q, se obtuvo %q", expectedContentType, actual)
				}

				// Decodificar el JSON de respuesta
				var actualBody healthResponse
				err = json.NewDecoder(rr.Body).Decode(&actualBody)
				if err != nil {
					t.Fatalf("No se pudo decodificar el JSON de respuesta: %v", err)
				}

				// Verificar el valor del campo "status"
				if actualBody.Status != tc.expectedBody.Status {
					t.Errorf("Status incorrecto: se esperaba %q, se obtuvo %q", tc.expectedBody.Status, actualBody.Status)
				}

				// Verificar el valor del campo "service" (mapeado desde la etiqueta JSON 'service')
				if actualBody.Services != tc.expectedBody.Services {
					t.Errorf("Service incorrecto: se esperaba %q, se obtuvo %q", tc.expectedBody.Services, actualBody.Services)
				}
			}
		})
	}
}
