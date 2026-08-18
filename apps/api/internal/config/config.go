package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	DatabaseURL     string
}

func Load() Config {
	readTimeout, err := getEnvAsDuration("READ_TIMEOUT_SECS", 5)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	writeTimeout, err := getEnvAsDuration("WRITE_TIMEOUT_SECS", 10)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	return Config{
		Port:         getEnv("PORT", "8080"),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsDuration(key string, fallbackSeconds int) (time.Duration, error) {
	// 1. Buscamos la variable en el entorno de forma nativa
	valueStr, exists := os.LookupEnv(key)

	// CASO 1: Variable ausente en el entorno -> usar fallback de forma silenciosa
	if !exists {
		return time.Duration(fallbackSeconds) * time.Second, nil
	}

	// CASO 2: Variable presente pero está vacía (ej: PORT="") -> también usamos fallback
	if valueStr == "" {
		return time.Duration(fallbackSeconds) * time.Second, nil
	}

	// Intentamos parsear el string a entero
	seconds, err := strconv.Atoi(valueStr)
	if err != nil {
		// CASO 3: Variable presente pero es inválida (ej: READ_TIMEOUT_SECS="diez")
		// Devolvemos una duración en cero (o el fallback) acompañado del error real
		return 0, fmt.Errorf("invalid value for environment variable %s: %w", key, err)
	}

	// CASO 4: Variable presente y válida -> usar su valor parseado
	return time.Duration(seconds) * time.Second, nil
}
