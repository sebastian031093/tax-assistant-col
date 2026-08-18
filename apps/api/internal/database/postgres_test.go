package database_test

import (
	"context"
	"strings"
	"tax-assistant-col/internal/database"
	"testing"
)

func TestNewPostgresPool_InfrastructureErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("Fails cleanly when PostgreSQL URL is structurally invalid", func(t *testing.T) {
		// URLs with unparseable control characters or bad schemes trigger ParseConfig errors
		invalidURL := "postgres://user:pass@[:123456]/database"

		pool, err := database.NewPostgresPool(ctx, invalidURL)

		if err == nil {
			t.Fatal("expected an error due to invalid URL structure, got nil")
		}
		if pool != nil {
			t.Error("expected pool instance pointer to be nil on failure paths")
		}
		if !strings.Contains(err.Error(), "error al parsear la URL") {
			t.Errorf("expected context error message wrapper, got: %v", err)
		}
	})

	t.Run("Fails and respects context cancellation during connection Ping phase", func(t *testing.T) {
		// Create a context that is pre-cancelled before passing it to the pool builder
		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel instantly

		// Use a valid URL format so it passes the config parsing step, but hits the ping
		validURLStr := "postgres://postgres:secret@localhost:5432/tax_db"

		pool, err := database.NewPostgresPool(canceledCtx, validURLStr)

		if err == nil {
			t.Fatal("expected connection to fail due to a canceled context context, got nil")
		}
		if pool != nil {
			t.Error("expected pool to be nil or completely destroyed on context cancellation")
		}

		// The error should mention that the context wrapper block aborted or failed the connection check
		if !strings.Contains(err.Error(), "falló la verificación de conexión") {
			t.Errorf("expected custom error wrapping for failed verification, got: %v", err)
		}
	})
}
