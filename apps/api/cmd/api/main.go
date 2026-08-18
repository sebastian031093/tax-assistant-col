package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tax-assistant-col/internal/config"
	"tax-assistant-col/internal/database" // Ajusta según la ruta real de tu paquete
	"tax-assistant-col/internal/server"
)

func main() {
	// Crear el contexto raíz que escucha las señales de apagado del sistema (Ctrl+C, SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. Load config
	cfg := config.Load()

	// Validar que la variable no venga vacía antes de intentar conectar
	if cfg.DatabaseURL == "" {
		fmt.Fprintln(os.Stderr, "Error crítico: la variable de entorno DATABASE_URL no está configurada.")
		os.Exit(1)
	}

	// 2 y 3. Create PostgreSQL pool & Verify connection
	log.Println("Inicializando y verificando conexión con PostgreSQL...")
	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		// Falle con un mensaje claro cuando DATABASE_URL sea incorrecta o el servidor no esté listo
		fmt.Fprintf(os.Stderr, "Fallo en la inicialización de la infraestructura: %v\n", err)
		os.Exit(1)
	}

	// 6. Close pool (Se ejecutará de último de forma segura al apagar la app)
	defer func() {
		log.Println("Cerrando el pool de conexiones de PostgreSQL...")
		pool.Close()
	}()

	log.Println("Conexión a PostgreSQL verificada exitosamente.")

	// 4 y 5. Inject pool into server & Run server
	// Nota: Pasamos el 'ctx' de ciclo de vida y el 'pool' al servidor
	if err := server.Run(ctx, cfg, pool); err != nil {
		fmt.Fprintf(os.Stderr, "Error en la ejecución del servidor: %v\n", err)
		os.Exit(1)
	}

	log.Println("Aplicación finalizada correctamente.")
}
