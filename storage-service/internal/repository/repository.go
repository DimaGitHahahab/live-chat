package repository

import (
	"context"

	"storage-service/internal/domain"
	"storage-service/internal/repository/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	AddMessage(context.Context, *domain.Message) (*domain.Message, error)
	GetMessage(context.Context, int) (*domain.Message, error)
}

type postgresRepo struct {
	*queries.SqlQueries
	pool *pgxpool.Pool
}

func New(pgxPool *pgxpool.Pool) Repository {
	return &postgresRepo{
		SqlQueries: queries.NewPostgresQueries(pgxPool),
		pool:       pgxPool,
	}
}
