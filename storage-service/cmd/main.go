package main

import (
	"storage-service/internal/app"
	"storage-service/pkg/config"
	"storage-service/pkg/logger"
)

func main() {
	logger.Init()

	cfg := config.LoadConfig()

	a := app.New(cfg)

	a.Run()
}
