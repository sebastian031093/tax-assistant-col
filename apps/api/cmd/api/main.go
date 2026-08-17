package main

import (
	"log"
	"tax-assistant-col/internal/config"
	"tax-assistant-col/internal/server"
)

func main() {
	// 1. Carga la configuración
	cfg := config.Load()

	if err := server.Run(cfg); err != nil {
		log.Fatal(err)
	}

	// 2. Inicia la aplicación y el servidor
	server.Run(cfg)
}
