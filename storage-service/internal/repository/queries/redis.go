package queries

import (
	"context"
	"encoding/json"
	"errors"

	"storage-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

const key = "messages"

type RedisQueries struct {
	redis *redis.Client
}

func NewRedisQueries(redis *redis.Client) *RedisQueries {
	return &RedisQueries{
		redis: redis,
	}
}

func (r RedisQueries) AddMessage(ctx context.Context, message *domain.Message) error {
	encodedMsg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.redis.LPush(ctx, key, encodedMsg).Err()
}

func (r RedisQueries) GetNumberOfMessages(ctx context.Context) (int, error) {
	num, err := r.redis.LLen(ctx, key).Result()
	return int(num), err
}

func (r RedisQueries) GetMessages(ctx context.Context) ([]domain.Message, error) {
	encodedMessages, err := r.redis.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	reverse(encodedMessages)

	return decodeMessages(encodedMessages)
}

func reverse(m []string) {
	for i := 0; i < len(m)/2; i++ {
		m[i], m[len(m)-i-1] = m[len(m)-i-1], m[i]
	}
}

func decodeMessages(encodedMessages []string) ([]domain.Message, error) {
	messages := make([]domain.Message, len(encodedMessages))
	for i, e := range encodedMessages {
		err := json.Unmarshal([]byte(e), &messages[i])
		if err != nil {
			return nil, err
		}
	}

	return messages, nil
}

func (r RedisQueries) RemoveOldestMessage(ctx context.Context) error {
	if _, err := r.redis.RPop(ctx, key).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
	}

	return nil
}
