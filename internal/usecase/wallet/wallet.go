package wallet

import (
	"context"
	errors2 "errors"
	"math/big"
	"task/internal/entity"
	"task/pkg/api/filter"
	"task/pkg/utils/map_converter"
)

func (uc *Usecase) GetAllMusicWithCategory(c context.Context, coin string) error {

	err := uc.wallet.CreateWallet(c, coin)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) EditBalance(c context.Context, trans *entity.Transactional) error {
	amountBalance, err := uc.wallet.GetWallet(c, &filter.Options{
		Fields: []filter.Field{
			{
				Name:     "id",
				Operator: filter.OperatorEq,
				Value:    trans.ID,
			},
		},
	})
	if err != nil {
		return err
	}

	res, err := editBalance(amountBalance.Money, trans.Amount, trans.TransactionalType)
	if err != nil {
		return err
	}

	m := map_converter.ConvertStructToMap(&entity.TransactionalDTO{
		WalletID:          amountBalance.ID,
		TransactionalType: trans.TransactionalType,
		Amount:            trans.Amount.String(),
		UpdatedBalance:    res.String(),
	})

	mWallet := map_converter.ConvertStructToMap(&entity.WalletDB{
		ID:    amountBalance.ID,
		Coin:  amountBalance.Coin,
		Money: res.String(),
	})

	err = uc.wallet.IncOrDecrBalance(c, m, amountBalance.ID, mWallet)
	if err != nil {
		return err
	}

	return nil
}

func editBalance(amountBalance string, amount *big.Int, trans string) (*big.Int, error) {
	balance := new(big.Int)
	balance.SetString(amountBalance, 10)

	final := big.Int{}

	if trans == "1" {
		final.Add(balance, amount)
		return &final, nil
	} else {
		diff := final.Sub(balance, amount)
		zero := new(big.Int)

		res := diff.Cmp(zero)
		if res == -1 {
			return nil, errors2.New("insufficient funds")
		}

		return diff, nil
	}
}

func (uc *Usecase) MoneyTransactional(ctx context.Context, transactional *entity.MoneyTransactional) error {
	const decreaseMoney = "0"
	const increaseMoney = "1"
	amountBalanceFrom, err := uc.wallet.GetWallet(ctx, &filter.Options{
		Fields: []filter.Field{
			{
				Name:     "id",
				Operator: filter.OperatorEq,
				Value:    transactional.FromWalletId,
			},
		},
	})
	if err != nil {
		return err
	}

	balanceFrom, err := editBalance(amountBalanceFrom.Money, transactional.Value, decreaseMoney)
	if err != nil {
		return err
	}

	mBalanceFrom := map_converter.ConvertStructToMap(&entity.WalletDB{
		ID:    amountBalanceFrom.ID,
		Coin:  amountBalanceFrom.Coin,
		Money: balanceFrom.String(),
	})

	mTransFrom := map_converter.ConvertStructToMap(&entity.TransactionalDTO{
		WalletID:          amountBalanceFrom.ID,
		TransactionalType: decreaseMoney,
		Amount:            transactional.Value.String(),
		UpdatedBalance:    balanceFrom.String(),
	})

	amountBalanceTo, err := uc.wallet.GetWallet(ctx, &filter.Options{
		Fields: []filter.Field{
			{
				Name:     "id",
				Operator: filter.OperatorEq,
				Value:    transactional.ToWalletId,
			},
		},
	})
	if err != nil {
		return err
	}

	balanceTo, err := editBalance(amountBalanceTo.Money, transactional.Value, increaseMoney)
	if err != nil {
		return err
	}

	mBalanceTo := map_converter.ConvertStructToMap(&entity.WalletDB{
		ID:    amountBalanceTo.ID,
		Coin:  amountBalanceTo.Coin,
		Money: balanceTo.String(),
	})

	mTransTo := map_converter.ConvertStructToMap(&entity.TransactionalDTO{
		WalletID:          amountBalanceTo.ID,
		TransactionalType: increaseMoney,
		Amount:            transactional.Value.String(),
		UpdatedBalance:    balanceTo.String(),
	})

	err = uc.wallet.MoneyTransactional(ctx, mBalanceFrom, mTransFrom, mBalanceTo, mTransTo)
	if err != nil {
		return err
	}

	return nil
}
