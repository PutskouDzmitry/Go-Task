package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"task/internal/config"

	"go.uber.org/zap"
)

type postgresqlClient struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewPostgresqlClient is a function that creates a new PostgreSQL storage client.
//
// It takes in a pointer to a postgresqlConfig struct and a logger.
// Returns a Storage interfaces and an error.
func NewPostgresqlClient(cfg *config.Config, logger *zap.Logger) (Storage, error) {
	pgConfig, err := newPostgresqlConfig(cfg)
	if err != nil {
		logger.Error("failed to get postgresql poolConfig, err", zap.Error(err))
		return nil, err
	}

	poolConfig, err := pgConfig.getPostgresqlConfig(logger)
	if err != nil {
		logger.Error("failed to get postgresql poolConfig, err", zap.Error(err))
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Error("failed to create postgresql pool, err", zap.Error(err))
		return nil, err
	}
	return &postgresqlClient{pool: pool, logger: logger}, nil
}

// Ping pings the PostgreSQL client.
//
// ctx: The context for the function.
// Returns an error if there was a problem pinging the client.
func (db *postgresqlClient) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Close closes the PostgreSQL client connection.
func (db *postgresqlClient) Close() {
	db.pool.Close()
}

// Acquire acquires a connection from the PostgreSQL client pool.
//
// ctx - The context.Context object.
// Return type - *pgxpool.Conn, error.
func (db *postgresqlClient) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return db.pool.Acquire(ctx)
}

// Begin starts a new transaction in the PostgreSQL client.
//
// ctx - The context to use for the transaction.
// Returns a pgx.Tx, which represents a PostgreSQL transaction, and an error if any occurred.
func (db *postgresqlClient) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.pool.Begin(ctx)
}

// BeginTx begins a transaction on the PostgreSQL client.
//
// ctx: the context.Context object.
// txOptions: the pgx.TxOptions object.
// Returns a pgx.Tx object and an error.
func (db *postgresqlClient) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return db.pool.BeginTx(ctx, txOptions)
}

// GetTxOptionRead returns the reading transaction options for reading from the PostgreSQL client.
//
// No parameters.
// Returns a pgx.TxOptions struct.
func (db *postgresqlClient) GetTxOptionRead() pgx.TxOptions {
	return pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	}
}

// GetTxOptionWrite returns the transaction options for performing a write operation in the PostgreSQL client.
//
// This function does not take any parameters.
// It returns a pgx.TxOptions struct.
func (db *postgresqlClient) GetTxOptionWrite() pgx.TxOptions {
	return pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}
}

// QueryRow executes a query that is expected to return at most one row.
//
// ctx is the context.Context.
// query is the SQL query to execute.
// args are the arguments to replace placeholders in the query.
// Returns a pgx.Row that represents the result of the query.
func (db *postgresqlClient) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, query, args...)
}

// Query executes a SQL query on the PostgreSQL database.
//
// ctx: The context.Context to use for the query.
// query: The SQL query to execute.
// args: The optional arguments to pass to the query.
// Returns: The pgx.Rows result of the query and an error if any.
func (db *postgresqlClient) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(ctx, query, args...)
}

// Exec executes a query on the PostgreSQL database.
//
// ctx - the context to use for the operation.
// query - the query to execute.
// args - optional arguments to be used in the query.
// Returns the command tag and an error if any.
func (db *postgresqlClient) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, query, args...)
}
