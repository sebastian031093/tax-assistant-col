package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tax-assistant-col/internal/config"
	"tax-assistant-col/internal/database"
	"tax-assistant-col/internal/server"
)

func main() {
	// main is purely responsible for handling the OS exit code status
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "application failed: %v\n", err)
		os.Exit(1)
	}
	log.Println("Application finished cleanly.")
}

// run contains the entire bootstrap logic, allowing all defers to execute safely
func run() error {
	// Create the root context listening for OS signals (Ctrl+C, SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. Load config
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("the environment variable DATABASE_URL is not configured")
	}

	// 2 & 3. Create PostgreSQL pool & Verify connection
	log.Println("Initializing and verifying connection with PostgreSQL...")
	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("infrastructure initialization failed: %w", err)
	}

	// 6. Close pool (Guaranteed to execute when run() returns an error OR exits successfully)
	defer func() {
		log.Println("Closing PostgreSQL connection pool...")
		pool.Close()
	}()

	log.Println("PostgreSQL connection verified successfully.")

	// 4 & 5. Inject pool into server & Run server
	if err := server.Run(ctx, cfg, pool); err != nil {
		return fmt.Errorf("server execution failed: %w", err)
	}

	return nil
}
