package validate

import "testing"

func TestIsPhoneNumber(t *testing.T) {
	tests := []struct {
		phone    string
		expected bool
	}{
		{"+79108184614", true},
		{"+79123456789", true},
		{"89012345678", true},
		{"1234567890", false},
		{"+1234567890", false},
	}

	for _, test := range tests {
		result := IsPhoneNumber(test.phone)
		if result != test.expected {
			t.Errorf("Expected IsPhoneNumber(%s) to be %v, but got %v", test.phone, test.expected, result)
		}
	}
}

func TestIsEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"test.email@example.co.uk", true},
		{"test123", false},
		{"test@", false},
		{"@example.com", false},
	}

	for _, test := range tests {
		result := IsEmail(test.email)
		if result != test.expected {
			t.Errorf("Expected IsEmail(%s) to be %v, but got %v", test.email, test.expected, result)
		}
	}
}
