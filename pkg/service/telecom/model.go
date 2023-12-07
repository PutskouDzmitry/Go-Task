package telecom

import (
	"fmt"
	"math"
	"math/rand"
)

type SmsService interface {
	Auth() (bool, error)
	SendCode(phone string) (string, error)
	CallCode(phone string) (string, error)
}

func generateCode(len int) string {
	if len < 1 {
		panic("len must be created than 0")
	}
	code := rand.Intn(int(math.Pow10(len))-int(math.Pow10(len-1))) + int(math.Pow10(len-1))
	return fmt.Sprint(code)
}
