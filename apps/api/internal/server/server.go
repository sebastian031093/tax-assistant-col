package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/internal/config" // Reemplaza "mi-api" por el nombre de tu módulo en go.mod
)

// Run inicializa y arranca el servidor gestionando el apagado controlado
func Run(cfg config.Config) {
	mux := newRouter()

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  15 * time.Second,
	}

	// 1. Escuchar señales del sistema operativo usando un contexto cancelable
	// Capturará Ctrl+C (SIGINT) o terminación de la nube (SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Ejecutar ListenAndServe en una goroutine (hilo asíncrono)
	// Hacemos esto porque ListenAndServe bloquea el flujo principal. Al delegarlo,
	// permitimos que el código de abajo siga corriendo y espere las señales.
	go func() {
		log.Printf("Starting server on port :%s...\n", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed down: %v", err)
		}
	}()

	// El canal se bloquea aquí hasta que el OS envíe un Ctrl+C o SIGTERM
	<-ctx.Done()
	log.Println("Shutting down server gracefully...")

	// 3. Crear una ventana de tiempo límite (timeout) para el apagado seguro
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Invocamos Shutdown. Esperará peticiones activas un máximo de 10 segundos.
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting clean and tidy.")
}
