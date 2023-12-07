package wallet

import (
	errors2 "errors"
	"github.com/gin-gonic/gin"
	pkgError "github.com/pkg/errors"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"task/internal/entity"
	"task/pkg/api/errors"
)

// CreateWallet
// @Summary Добавить кошелек
// @Tags wallet
// @Accept json
// @Param coin path string true "Coin"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /wallet/create/{coin} [get]
func (h *Handler) CreateWallet(c *gin.Context) {
	typeOfCoin := c.Param("coin")

	if typeOfCoin == "" {
		errors.AbortWithBadRequest(c, pkgError.New("param is empty"))
		return
	}

	if checkCoin(typeOfCoin) != true {
		errors.AbortWithBadRequest(c, pkgError.New("the system doesn't suppose your type of coin"))
		return
	}

	err := h.wallet.GetAllMusicWithCategory(c, typeOfCoin)
	if err != nil {
		errors.NewErrorResponse(c, err)
		return
	}

	c.Status(http.StatusOK)

}

func checkCoin(coin string) bool {
	allCoins := []string{
		"eth",
		"bitcoin",
		"ethereum",
		"litecoin",
	}

	for _, value := range allCoins {
		if strings.ToLower(value) == coin {
			return true
		}
	}

	return false
}

func (h *Handler) EditBalance(c *gin.Context) {

	var dto entity.TransactionalReq

	if err := c.ShouldBindJSON(&dto); err != nil {
		errors.AbortWithBadRequest(c, err)
		return
	}

	amount, err := changeAmount(dto.Amount)
	if err != nil {
		errors.AbortWithBadRequest(c, err)
		return
	}

	if !(dto.TransactionalType == "0" || dto.TransactionalType == "1") {
		errors.AbortWithBadRequest(c, errors2.New("incorrect TransactionalType"))
		return
	}

	err = h.wallet.EditBalance(c, &entity.Transactional{
		ID:                dto.ID,
		Amount:            amount,
		TransactionalType: dto.TransactionalType,
	})
	if err != nil {
		errors.AbortWithBadRequest(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func changeAmount(amount string) (*big.Int, error) {
	str := strings.Split(amount, ".")
	if len(str) == 1 {
		if _, err := strconv.Atoi(str[0]); err != nil {
			return nil, errors2.New("your number isn't number")
		}
		final := new(big.Int)

		base := new(big.Int)
		base.SetString(amount, 10)

		forMulBase := new(big.Int)
		forMulBase.SetString("1000000000000000000", 10)

		final.Mul(forMulBase, base)
		return final, nil
	} else {
		if len(str) != 2 {
			return nil, errors2.New("your amount isn't correct")
		}
		if _, err := strconv.Atoi(str[1]); err != nil {
			return nil, errors2.New("your number isn't number")
		}
		base := new(big.Int)
		forBase := new(big.Int)

		fraction := new(big.Int)
		forFraction := new(big.Int)

		final := new(big.Int)

		forMulBase := new(big.Int)
		forMulBase.SetString("1000000000000000000", 10)
		base.SetString(str[0], 10)
		forBase.Mul(base, forMulBase)

		strFraction := generateFraction("1", 0, 18-len(str[1]))
		forMulFraction := new(big.Int)
		forMulFraction.SetString(strFraction, 10)
		fraction.SetString(str[1], 10)
		forFraction.Mul(forMulFraction, fraction)

		final.Add(forBase, forFraction)
		return final, nil
	}
}

func generateFraction(str string, n, counter int) string {
	if n != counter {
		str = str + "0"
		n++
		return generateFraction(str, n, counter)
	}
	return str
}

func (h *Handler) MoneyTransactional(c *gin.Context) {
	var dto entity.MoneyTransactionalReq

	if err := c.ShouldBindJSON(&dto); err != nil {
		errors.AbortWithBadRequest(c, err)
		return
	}

	amount, err := changeAmount(dto.Value)
	if err != nil {
		errors.AbortWithBadRequest(c, err)
		return
	}

	err = h.wallet.MoneyTransactional(c, &entity.MoneyTransactional{
		FromWalletId: dto.FromWalletId,
		ToWalletId:   dto.ToWalletId,
		Value:        amount,
	})
	if err != nil {
		errors.AbortWithBadRequest(c, err)
		return
	}

}
