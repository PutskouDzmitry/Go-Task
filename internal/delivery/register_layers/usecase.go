package register_layers

import (
	"task/internal/usecase/wallet"
)

type GUsecase struct {
	wallet *wallet.Usecase
}

func NewGUsecase(repo *GRepository) *GUsecase {
	return &GUsecase{
		wallet: wallet.NewUsecase(repo.wallet),
	}
}
