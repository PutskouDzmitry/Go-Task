package validate

import "regexp"

var regex *regexp.Regexp

func IsPhoneNumber(phone string) bool {
	regex = regexp.MustCompile(`^(\+7|8)\d{10}$`)
	return regex.MatchString(phone)
}

func IsEmail(email string) bool {
	regex = regexp.MustCompile(`[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)*@[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)*`)
	return regex.MatchString(email)
}
