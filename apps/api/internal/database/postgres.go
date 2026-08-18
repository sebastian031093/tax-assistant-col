package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool inicializa y verifica un pool de conexiones a PostgreSQL.
func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	// 1. Validar o parsear la URL
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error al parsear la URL de la base de datos: %w", err)
	}

	// 2. Crear el pool (no abre conexiones físicas inmediatamente)
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error al crear el pool de conexiones: %w", err)
	}

	// Crear un contexto con timeout de 5 segundos exclusivamente para el Ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 3. Ejecutar Ping con un timeout para verificar la conexión real
	if err := pool.Ping(pingCtx); err != nil {
		// 5. Cerrar el pool si la inicialización falla (evita fugas de memoria)
		pool.Close()
		// 4. Devolver un error con contexto si falla
		return nil, fmt.Errorf("falló la verificación de conexión (Ping) a PostgreSQL: %w", err)
	}

	return pool, nil
}
