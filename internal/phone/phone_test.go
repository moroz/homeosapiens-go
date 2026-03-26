package phone_test

import (
	"strings"
	"testing"

	"github.com/moroz/homeosapiens-go/internal/countries"
	"github.com/moroz/homeosapiens-go/internal/phone"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeNumber(t *testing.T) {
	examples := []struct {
		input, expected string
		valid           bool
	}{
		// Different forms of numbers within Poland
		{"+48 555-123-456", "+48555123456", true},
		{"(+48) 555-123-456", "+48555123456", true},
		{"+48555123456", "+48555123456", true},
		{"555-123-456", "+48555123456", true},
		{"+49 171 234 5678", "+491712345678", true},
		{"555123456", "+48555123456", true},
		{"0048555123456", "+48555123456", true},
		{"48555123456", "+48555123456", true},

		// Obviously invalid numbers
		{"", "", false},
		{"abc", "", false},
		{"123", "", false},           // too short
		{"+999555123456", "", false}, // invalid country code
		{"++48555123456", "", false}, // malformed prefix
	}

	for _, example := range examples {
		actual, err := phone.NormalizeNumber(example.input, "")
		if example.valid {
			assert.NoError(t, err)
		}
		assert.Equal(t, example.expected, actual)
	}
}

func TestExamplePhoneNumber(t *testing.T) {
	for _, iso := range countries.EuMemberStatesISO {
		number := phone.ExamplePhoneNumber(iso)
		assert.NotEmpty(t, number)

		expected := strings.ReplaceAll(number, " ", "")
		actual, err := phone.NormalizeNumber(number, "")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}
