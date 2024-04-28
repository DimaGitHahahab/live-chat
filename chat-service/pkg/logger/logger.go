package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{})
}
