package main

import (
	"chat-service/internal/app"
	"chat-service/pkg/config"
	"chat-service/pkg/logger"
)

func main() {
	logger.Init()

	cfg := config.LoadConfig()

	a := app.New(cfg)

	a.Run()
}
