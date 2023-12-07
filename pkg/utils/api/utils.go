package api

import (
	"MommyCO/pkg/api/jwt"
	"fmt"
	j "github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

func GetMapClaims(c *gin.Context, jwt *jwt.JWT) (j.MapClaims, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, fmt.Errorf("authorization header required")
	}

	mapClaims, err := jwt.Parse(token)
	if err != nil {
		return nil, err
	}

	return mapClaims, nil
}
