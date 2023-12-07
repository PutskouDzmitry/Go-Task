package wallet

import (
	"context"
	"github.com/gin-gonic/gin"
	"task/internal/entity"
)

type Delivery interface {
	GetAllMusicWithCategory(c context.Context, coin string) error
	EditBalance(c context.Context, trans *entity.Transactional) error
	MoneyTransactional(ctx context.Context, transactional *entity.MoneyTransactional) error
}

type Handler struct {
	wallet Delivery
}

func NewHandler(wallet Delivery) *Handler {
	return &Handler{
		wallet: wallet,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/wallet")

	group.GET("/create/:coin", h.CreateWallet)
	group.POST("/", h.EditBalance)
	group.POST("/trans", h.MoneyTransactional)
}
