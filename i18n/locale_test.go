package i18n_test

import (
	"testing"

	"github.com/moroz/homeosapiens-go/i18n"
	"github.com/stretchr/testify/assert"
)

func TestResolveLocaleFromAcceptLanguageHeader(t *testing.T) {
	examples := []struct {
		header   string
		expected string
	}{
		{"en-US,en;q=0.9", "en"},
		{"en", "en"},
		{"en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7", "en"},
		{"en-IN", "en"},
		{"en-IN,en;q=0.9", "en"},
		{"hi-IN,hi;q=0.9,en-IN;q=0.8,en;q=0.7", "en"},
		{"pl-PL", "pl"},
		{"pl-PL,pl;q=0.9,en;q=0.8", "pl"},
		{"zh-CN", ""},
		{"", ""},
	}

	for _, example := range examples {
		actual, err := i18n.ResolveLocaleFromAcceptLanguageHeader(example.header)
		assert.NoError(t, err)
		assert.Equal(t, example.expected, actual)
	}
}

func TestResolveLocale(t *testing.T) {
	examples := []struct {
		langParam string
		header    string
		expected  string
	}{
		{"", "en-US,en;q=0.9", "en"},
		{"pl", "en-US,en;q=0.9", "pl"},
		{"en", "en-US,en;q=0.9", "en"},
		{"zh", "en-US,en;q=0.9", "en"},
		{"de", "en-US,en;q=0.9", "en"},
		{"", "pl-PL", "pl"},
		{"en", "pl-PL", "en"},
	}

	for _, example := range examples {
		actual := i18n.ResolveLocale(example.langParam, example.header)
		assert.Equal(t, example.expected, actual)
	}
}
