package entity

import (
	"fmt"
	"math/big"
)

type Wallet struct {
	ID    string   `json:"id"`
	Coin  string   `json:"coin"`
	Money *big.Int `json:"money"`
}

type WalletDB struct {
	ID    string `json:"id"`
	Coin  string `json:"coin"`
	Money string `json:"money"`
}

func (c *WalletDB) Columns(alias string) []string {
	return []string{
		fmt.Sprintf("%s.id", alias),
		fmt.Sprintf("%s.coin", alias),
		fmt.Sprintf("%s.money", alias),
	}
}

func (c *WalletDB) Dest() []interface{} {
	return []interface{}{
		&c.ID,
		&c.Coin,
		&c.Money,
	}
}

type TransactionalReq struct {
	ID                string `json:"wallet_id" binding:"required"`
	Amount            string `json:"amount" binding:"required"`
	TransactionalType string `json:"transaction_type" binding:"required"`
}

type Transactional struct {
	ID                string   `json:"wallet_id" binding:"required"`
	Amount            *big.Int `json:"amount" binding:"required"`
	TransactionalType string   `json:"transaction_type" binding:"required"`
}

type TransactionalDTO struct {
	WalletID          string `json:"wallet_id"`
	TransactionalType string `json:"transactional_type"`
	Amount            string `json:"amount"`
	UpdatedBalance    string `json:"updated_balance"`
}

func (c *TransactionalDTO) Columns(alias string) []string {
	return []string{
		fmt.Sprintf("%s.id", alias),
		fmt.Sprintf("%s.wallet_id", alias),
		fmt.Sprintf("%s.transactional_type", alias),
		fmt.Sprintf("%s.amount", alias),
		fmt.Sprintf("%s.updated_balance", alias),
	}
}

func (c *TransactionalDTO) Dest() []interface{} {
	return []interface{}{
		&c.WalletID,
		&c.UpdatedBalance,
		&c.Amount,
		&c.TransactionalType,
	}
}

type MoneyTransactionalReq struct {
	FromWalletId string `json:"from_wallet_id"`
	ToWalletId   string `json:"to_wallet_id"`
	Value        string `json:"value"`
}

type MoneyTransactional struct {
	FromWalletId string   `json:"from_wallet_id"`
	ToWalletId   string   `json:"to_wallet_id"`
	Value        *big.Int `json:"value"`
}
