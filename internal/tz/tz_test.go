package tz_test

import (
	"testing"

	"github.com/moroz/homeosapiens-go/internal/tz"
	"github.com/stretchr/testify/assert"
)

func TestGuessRegionByTimezone(t *testing.T) {
	examples := []struct {
		iso    string
		tzname string
	}{
		{"TW", "Asia/Taipei"},
		{"IN", "Asia/Kolkata"},
		{"PL", "Europe/Warsaw"},
		{"CN", "Asia/Shanghai"},
		{"HK", "Asia/Hong_Kong"},
		{"MO", "Asia/Macau"},
		{"US", "America/New_York"},
		{"ZA", "Africa/Johannesburg"},
	}

	for _, example := range examples {
		actual := tz.GuessRegionByTimezone(example.tzname)
		assert.Equal(t, example.iso, actual.IsoCode)
	}
}
