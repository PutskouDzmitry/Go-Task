package register_layers

import (
	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"
	"task/internal/interfaces"
	"task/internal/repository/wallet"
	"task/pkg/database/postgresql"
)

type GRepository struct {
	wallet interfaces.Wallet
}

func NewGRepository(
	db postgresql.Storage,
	log *zap.Logger,
) *GRepository {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &GRepository{
		wallet: wallet.NewWallet(db, queryBuilder, log),
	}
}
