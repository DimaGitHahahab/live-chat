package adapter

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

type MessageHandler func(message *sarama.ConsumerMessage) error

type Consumer struct {
	handler MessageHandler
	brokers []string
	topic   string
	group   string
}

func NewConsumer(handler MessageHandler, brokers []string, topic, group string) *Consumer {
	return &Consumer{
		handler: handler,
		brokers: brokers,
		topic:   topic,
		group:   group,
	}
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Infoln("Consumer message channel was closed")
				return nil
			}

			err := c.handler(message)
			if err != nil {
				log.Errorln("Error while handling consumer message: ", err)
				return nil
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			session.Commit()
			return nil
		}
	}
}

func (c *Consumer) StartConsuming(ctx context.Context) error {
	consumerGroup, err := sarama.NewConsumerGroup(c.brokers, c.group, getConsumerConfig())
	if err != nil {
		return err
	}
	log.Infoln("Consumer group was created")
	defer consumerGroup.Close()

	// go func() {
	for {
		if err := consumerGroup.Consume(ctx, []string{c.topic}, c); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}
			log.Fatalln("Error from consumer group: ", err)
		}
		if ctx.Err() != nil {
			return err
		}
	}
	//}()
}

func getConsumerConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.DefaultVersion
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	return cfg
}
