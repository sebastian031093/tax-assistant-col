package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"tax-assistant-col/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Run configura, inicia y apaga de forma segura el servidor HTTP.
// Recibe el contexto raíz de la aplicación y el pool de conexiones inyectado.
func Run(ctx context.Context, cfg config.Config, pool *pgxpool.Pool) error {
	// TODO: Pasa el 'pool' a tu enrutador cuando vayas a configurar tus handlers/repositorios
	mux := newRouter()

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       15 * time.Second,
	}

	// 1. Canal con buffer para capturar errores de inicio del servidor HTTP
	serverErrors := make(chan error, 1)

	// 2. Ejecutar ListenAndServe de forma asíncrona
	go func() {
		log.Printf("Starting server on port :%s...\n", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// 3. Coordinar la espera de eventos usando el contexto inyectado desde main
	select {
	case <-ctx.Done():
		// Caso A: Se recibió una señal de apagado del sistema (Ctrl+C, SIGTERM) en main.go
		log.Println("Shutting down server gracefully...")

	case err := <-serverErrors:
		// Caso B: El servidor falló internamente antes de recibir una señal (ej: puerto ocupado)
		return fmt.Errorf("server failed to start: %w", err)
	}

	// 4. Ventana de tiempo límite (timeout) para el apagado seguro del servidor HTTP
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout*time.Second)
	defer cancel()

	// Intentar apagar el servidor HTTP de forma limpia liberando las conexiones activas
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exiting clean and tidy.")

	// Si el servidor se cerró de forma normal tras el Shutdown, retornamos nil
	return nil
}
