package app

import (
	"fmt"
	"net/url"
	"os"

	"chat-service/internal/service"
	"chat-service/pkg/config"
	"chat-service/pkg/signal"

	log "github.com/sirupsen/logrus"
)

type ClientApp struct {
	config  *config.Config
	sigQuit chan os.Signal
}

func NewClient(cfg *config.Config) *ClientApp {
	sigQuit := signal.GetShutdownChannel()

	return &ClientApp{
		config:  cfg,
		sigQuit: sigQuit,
	}
}

func (a *ClientApp) Run() {
	var c *service.Client

	go func() {
		log.Infoln("starting client...")
		var nickname string
		fmt.Println("Enter your nickname:")
		_, err := fmt.Scan(&nickname)
		if err != nil {
			log.Fatalln("failed to read nickname: ", err)
		}
		fmt.Println()

		c, err = service.NewClient(url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}, nickname)
		if err != nil {
			log.Fatalln("failed to create client: ", err)
		}

		c.Run()
	}()

	<-a.sigQuit
	log.Infoln("Gracefully shutting down client")

	if c != nil {
		if err := c.Stop(); err != nil {
			log.Fatalln("Failed to shutdown the client gracefully: ", err)
		}
	}

	log.Infoln("Client shutdown is successful")
}
