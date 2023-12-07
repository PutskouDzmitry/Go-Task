package validate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePhone(t *testing.T) {

	t.Run("prefix 8", func(t *testing.T) {
		validPhone := "89161234567"
		expected := "+79161234567"
		res, err := ValidatePhone(validPhone)
		if err != nil || res != expected {
			t.Errorf("Expected phone %s to be %s, got %s", validPhone, expected, res)
		}
	})

	t.Run("prefix 7", func(t *testing.T) {
		validPhone := "+79161234567"
		expected := "+79161234567"
		res, err := ValidatePhone(validPhone)
		if err != nil || res != expected {
			t.Errorf("Expected phone %s to be %s, got %s", validPhone, expected, res)
		}
	})

	t.Run("invalid phone", func(t *testing.T) {
		invalidPhone := "123"
		_, err := ValidatePhone(invalidPhone)
		if err == nil {
			assert.Error(t, err)
		}
	})
}

func TestValidateEmail(t *testing.T) {
	validEmail := "Test@Example.COM"
	expected := "test@example.com"
	res, err := ValidateEmail(validEmail)
	if err != nil || res != expected {
		t.Errorf("Expected email %s to be %s, got %s", validEmail, expected, res)
	}

	invalidEmail := "bad email"
	_, err = ValidateEmail(invalidEmail)
	if err == nil {
		t.Errorf("Expected email %s to be invalid", invalidEmail)
	}
}
