package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tax-assistant-col/internal/config"
	"time"
)

func Run(cfg config.Config) error {
	mux := newRouter()

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       15 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. Creamos un canal con buffer para capturar errores del servidor
	serverErrors := make(chan error, 1)

	// 2. Ejecutar ListenAndServe pasándole los errores al canal
	go func() {
		log.Printf("Starting server on port :%s...\n", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err // Enviamos el error al canal si el servidor falla al arrancar
		}
	}()

	// 3. AQUÍ VA TU BLOQUE SELECT: Coordina la espera de eventos
	select {
	case <-ctx.Done():
		// Caso A: Se recibió Ctrl+C o SIGTERM de la nube. Procedemos al apagado controlado.
		log.Println("Shutting down server gracefully...")

	case err := <-serverErrors:
		// Caso B: El servidor falló antes de recibir una señal (ej: puerto duplicado).
		log.Fatalf("Server failed to start: %v", err)
	}

	// 4. Ventana de tiempo límite (timeout) para el apagado seguro (solo aplica si entró al Caso A)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting clean and tidy.")

	return <-serverErrors
}
