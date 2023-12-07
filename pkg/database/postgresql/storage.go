package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Ping(ctx context.Context) error
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	GetTxOptionRead() pgx.TxOptions
	GetTxOptionWrite() pgx.TxOptions
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
}
