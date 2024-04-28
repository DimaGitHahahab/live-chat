package adapter

import (
	"time"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

type Producer struct {
	sarama.AsyncProducer
	TopicName string
}

func NewKafkaProducer(brokers []string, topic string) *Producer {
	config := getConfig()

	saramaProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama saramaProducer: ", err)
	}

	log.Infof("Sucessfully connected to Kafka")

	producer := &Producer{
		AsyncProducer: saramaProducer,
		TopicName:     topic,
	}

	go func() {
		for err := range saramaProducer.Errors() {
			log.Errorln("Failed to write access log entry: ", err)
		}
	}()

	return producer
}

func getConfig() *sarama.Config {
	c := sarama.NewConfig()
	c.Version = sarama.DefaultVersion
	c.Producer.RequiredAcks = sarama.WaitForLocal
	c.Producer.Compression = sarama.CompressionSnappy
	c.Producer.Flush.Frequency = 500 * time.Millisecond

	return c
}
