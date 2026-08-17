package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	// 1. Inicializamos el router interno real de la aplicación
	router := newRouter()

	// 2. Definición de la tabla de pruebas (Table-Driven Test)
	tests := []struct {
		name         string         // Nombre descriptivo del caso
		method       string         // Método HTTP a simular
		path         string         // Ruta de la petición
		wantStatus   int            // Código de estado HTTP esperado
		validateBody bool           // Bandera para validar el JSON de respuesta
		expectedBody healthResponse // Estructura JSON esperada en /health
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
			name:         "POST /health método no permitido",
			method:       "POST",
			path:         "/health",
			wantStatus:   http.StatusMethodNotAllowed, // 405 (Mux nativo maneja el prefijo "GET ")
			validateBody: false,
		},
		{
			name:         "GET /unknown ruta inexistente capturada por el 404 global",
			method:       "GET",
			path:         "/unknown",
			wantStatus:   http.StatusNotFound, // 404
			validateBody: false,
		},
	}

	// 3. Iteración sobre cada caso de la tabla
	for _, tc := range tests {
		// Ejecutamos cada caso como una sub-prueba independiente
		t.Run(tc.name, func(t *testing.T) {

			// Crear la petición HTTP falsa basada en el caso actual
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatalf("No se pudo crear la petición: %v", err)
			}

			// Crear el ResponseRecorder para capturar la respuesta en memoria
			rr := httptest.NewRecorder()

			// Enviamos la petición a través del método ServeHTTP del router real
			// Esto valida el enrutamiento completo y las restricciones de verbos HTTP
			router.ServeHTTP(rr, req)

			// Validar el código de estado obtenido
			if rr.Code != tc.wantStatus {
				t.Errorf("Código de estado incorrecto: se esperaba %d, se obtuvo %d", tc.wantStatus, rr.Code)
			}

			// Validar el contenido del JSON si el caso lo requiere
			if tc.validateBody {
				// Validar que el Header Content-Type sea application/json
				expectedContentType := "application/json"
				if actual := rr.Header().Get("Content-Type"); actual != expectedContentType {
					t.Errorf("Content-Type incorrecto: se esperaba %q, se obtuvo %q", expectedContentType, actual)
				}

				// Decodificar el cuerpo de la respuesta
				var actualBody healthResponse
				err = json.NewDecoder(rr.Body).Decode(&actualBody)
				if err != nil {
					t.Fatalf("No se pudo decodificar el JSON de respuesta: %v", err)
				}

				// Verificar el campo "status"
				if actualBody.Status != tc.expectedBody.Status {
					t.Errorf("Status incorrecto: se esperaba %q, se obtuvo %q", tc.expectedBody.Status, actualBody.Status)
				}

				// Verificar el campo "service"
				if actualBody.Services != tc.expectedBody.Services {
					t.Errorf("Services incorrecto: se esperaba %q, se obtuvo %q", tc.expectedBody.Services, actualBody.Services)
				}
			}
		})
	}
}
