package wallet

import (
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
	"task/pkg/database/postgresql"
)

type Wallet struct {
	db           postgresql.Storage
	queryBuilder sq.StatementBuilderType
	log          *zap.Logger
}

func NewWallet(db postgresql.Storage, queryBuilder sq.StatementBuilderType, log *zap.Logger) *Wallet {
	return &Wallet{db: db, queryBuilder: queryBuilder, log: log}
}
