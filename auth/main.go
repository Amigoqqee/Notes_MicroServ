package main

import (
	"auth/internal/config"
	"auth/internal/server"
	"log"
)

func main() {
	cfg := config.NewConfig()

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Printf("Ошибка при создании сервера: %v", err)
		return
	}

	if err := server.Serve(); err != nil {
		log.Printf("Ошибка запуска сервера: %v", err)
		return
	}
}
