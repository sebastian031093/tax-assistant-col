package main

import (
	"github.com/internal/config"
	"github.com/internal/server"
)

func main() {
	// 1. Carga la configuración
	cfg := config.Load()

	// 2. Inicia la aplicación y el servidor
	server.Run(cfg)
}
