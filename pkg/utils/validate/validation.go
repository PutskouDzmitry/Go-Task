package validate

import (
	"errors"
	"strings"
)

func ValidatePhone(phone string) (string, error) {
	if !IsPhoneNumber(phone) {
		return "", errors.New("invalid phone number")
	}

	if strings.HasPrefix(phone, "8") {
		return "+7" + phone[1:], nil
	}
	return phone, nil
}

func ValidateEmail(email string) (string, error) {
	if !IsEmail(email) {
		return "", errors.New("invalid email")
	}

	email = strings.ToLower(email)
	return email, nil
}
