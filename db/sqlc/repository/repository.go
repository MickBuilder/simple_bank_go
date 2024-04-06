package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewRepository(connPool *pgxpool.Pool) Repository {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
