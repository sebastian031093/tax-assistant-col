package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() Config {
	return Config{
		Port:         getEnv("PORT", "8080"),
		ReadTimeout:  getEnvAsDuration("READ_TIMEOUT_SECS", 5),
		WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT_SECS", 10),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsDuration(key string, fallbackSeconds int) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return time.Duration(fallbackSeconds) * time.Second
	}
	seconds, err := strconv.Atoi(valueStr)
	if err != nil {
		return time.Duration(fallbackSeconds) * time.Second
	}
	return time.Duration(seconds) * time.Second
}
