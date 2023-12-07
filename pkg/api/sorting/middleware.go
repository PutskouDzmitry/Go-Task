package sorting

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	SORT_OPTIONS_KEY = "sort_options"
	ASC_ORDER        = "ASC"
	DESC_ORDER       = "DESC"
)

func SortMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sortField := c.Query("sort_by")
		if sortField == "" {
			c.Next()
			return
		}

		fields := strings.Split(sortField, ",")

		sortOrder := c.Query("sort_order")
		if sortOrder == "" {
			c.Next()
			return
		}
		orders := strings.Split(sortOrder, ",")
		for _, order := range orders {
			if strings.ToUpper(order) != ASC_ORDER && strings.ToUpper(order) != DESC_ORDER {
				err := c.AbortWithError(http.StatusBadRequest, errors.New("bad sort order, must be asc/desc"))
				if err != nil {
					zap.L().Error("bad sort order, must be asc/desc")
				}
				return
			}
		}

		if len(orders) != len(fields) {
			err := c.AbortWithError(http.StatusBadRequest, errors.New("len of sort fields doesnt equal len of sort orders"))
			if err != nil {
				zap.L().Error("len of sort fields doesnt equal len of sort orders")
			}
		}

		var options Options
		options.Fields = fields
		options.Orders = orders
		c.Set(SORT_OPTIONS_KEY, options)
		c.Next()
	}
}

func GetOptionsFromContext(c *gin.Context) *Options {
	obj, ok := c.Get(SORT_OPTIONS_KEY)
	if !ok {
		return &Options{}
	}
	so := obj.(Options)
	return &so
}
