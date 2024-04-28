package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	HttpPort          string `envconfig:"HTTP_PORT" default:"8080"`
	KafkaBrokers      string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
	KafkaTopic        string `envconfig:"KAFKA_TOPIC" default:"chat.public.messages"`
	StorageServiceUrl string `envconfig:"STORAGE_URL" default:"localhost:8081"`
}

func LoadConfig() *Config {
	log.Infoln("Loading env variables")

	if err := godotenv.Load(); err != nil {
		log.Fatalln("Failed to load env variables: ", err)
	}

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalln("Failed to process env variables: ", err)
	}

	log.Infof("Loaded env variables successfully")
	return &cfg
}
