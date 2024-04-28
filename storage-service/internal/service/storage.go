package service

import (
	"context"
	"encoding/json"

	"storage-service/internal/adapter"
	"storage-service/internal/domain"
	"storage-service/internal/repository"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

const cacheSize = 10

type Storage struct {
	cache repository.Cache
	repo  repository.Repository
	ctx   context.Context
}

func NewStorage(ctx context.Context, repo repository.Repository, cache repository.Cache, kafkaBrokers []string, kafkaTopic string, group string) *Storage {
	s := &Storage{
		ctx:   ctx,
		cache: cache,
		repo:  repo,
	}

	consumer := adapter.NewConsumer(s.manageNewMessage, kafkaBrokers, kafkaTopic, group)

	go func() {
		if err := consumer.StartConsuming(ctx); err != nil {
			log.Fatalln("Failed to start consuming messages: ", err)
		}
	}()

	return s
}

func (s *Storage) manageNewMessage(saramaMessage *sarama.ConsumerMessage) error {
	var message domain.Message
	if err := json.Unmarshal(saramaMessage.Value, &message); err != nil {
		log.Errorln("Failed to unmarshal message from MQ: ", err)
	}

	savedMsg, err := s.repo.AddMessage(s.ctx, &message)
	if err != nil {
		return err
	}

	return s.updateCache(savedMsg)
}

func (s *Storage) updateCache(m *domain.Message) error {
	numberOfCachedMessages, err := s.cache.GetNumberOfMessages(s.ctx)
	if err != nil {
		return err
	}

	if numberOfCachedMessages >= cacheSize {
		err = s.cache.RemoveOldestMessage(s.ctx)
		if err != nil {
			return err
		}
	}

	return s.cache.AddMessage(s.ctx, m)
}

func (s *Storage) GetLastMessages() []domain.Message {
	msgs, err := s.cache.GetMessages(s.ctx)
	if err != nil {
		log.Errorln("Error while trying to get caches messages: ", err)
	}

	return msgs
}
