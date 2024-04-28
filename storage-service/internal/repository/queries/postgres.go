package queries

import (
	"context"

	"storage-service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlQueries struct {
	pool *pgxpool.Pool
}

func NewPostgresQueries(pgxPool *pgxpool.Pool) *SqlQueries {
	return &SqlQueries{pool: pgxPool}
}

const insertMessage = `INSERT INTO messages (sender, text, send_at) VALUES ($1, $2, $3) RETURNING id`

func (q *SqlQueries) AddMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	if err := q.pool.QueryRow(ctx, insertMessage, message.Sender, message.Text, message.SendAt).Scan(&message.Id); err != nil {
		return nil, err
	}

	return message, nil
}

const selectMessage = `SELECT id,sender, text, send_at  FROM messages WHERE id = $1`

func (q *SqlQueries) GetMessage(ctx context.Context, id int) (*domain.Message, error) {
	var message domain.Message
	if err := q.pool.QueryRow(ctx, selectMessage, id).Scan(&message.Id, &message.Sender, &message.Text, &message.SendAt); err != nil {
		return nil, err
	}

	return &message, nil
}
