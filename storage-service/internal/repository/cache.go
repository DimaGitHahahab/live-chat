package repository

import (
	"context"

	"storage-service/internal/domain"
	"storage-service/internal/repository/queries"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	AddMessage(ctx context.Context, message *domain.Message) error
	GetNumberOfMessages(ctx context.Context) (int, error)
	GetMessages(ctx context.Context) ([]domain.Message, error)
	RemoveOldestMessage(ctx context.Context) error
}

func NewCache(redis *redis.Client) Cache {
	return &redisCache{
		RedisQueries: queries.NewRedisQueries(redis),
	}
}

type redisCache struct {
	*queries.RedisQueries
}
