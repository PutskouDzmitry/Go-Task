package interfaces

import (
	"context"
	"task/internal/entity"
	"task/pkg/api/filter"
)

type Wallet interface {
	CreateWallet(ctx context.Context, coin string) error
	GetWallet(ctx context.Context, fo *filter.Options) (*entity.WalletDB, error)
	IncOrDecrBalance(ctx context.Context, m map[string]interface{}, id string, mWallet map[string]interface{}) error
	MoneyTransactional(ctx context.Context, mBalanceFrom, mTransFrom, mBalanceTo, mTransTo map[string]interface{}) error
}
