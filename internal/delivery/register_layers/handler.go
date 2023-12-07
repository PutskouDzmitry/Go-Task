package register_layers

import (
	"github.com/gin-gonic/gin"
	"task/internal/delivery/v1/wallet"
)

type GDelivery struct {
	wallet *wallet.Handler
}

func NewGDelivery(uc *GUsecase) *GDelivery {
	return &GDelivery{
		wallet: wallet.NewHandler(uc.wallet),
	}
}

func (h *GDelivery) RegisterRoutes(router *gin.RouterGroup) {
	h.wallet.RegisterRoutes(router)
}
