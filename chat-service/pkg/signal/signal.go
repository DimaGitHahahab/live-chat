package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func GetShutdownChannel() chan os.Signal {
	sigQuit := make(chan os.Signal)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	return sigQuit
}
