package wallet

import "task/internal/interfaces"

type Usecase struct {
	wallet interfaces.Wallet
}

func NewUsecase(wallet interfaces.Wallet) *Usecase {
	return &Usecase{
		wallet: wallet,
	}
}
