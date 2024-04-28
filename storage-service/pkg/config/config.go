package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	HttpPort string `envconfig:"HTTP_PORT" default:"8081"`

	KafkaBrokers string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
	KafkaTopic   string `envconfig:"KAFKA_TOPIC" default:"chat.public.messages"`
	KafkaGroup   string `envconfig:"KAFKA_GROUP" default:"1"`

	PostgresUrl string `envconfig:"PG_URL" default:"?????"`

	RedisAddr     string `envconfig:"REDIS_ADDR" required:"true"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" required:"true"`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`
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
