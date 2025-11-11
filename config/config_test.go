package config_test

import (
	"crypto/rand"
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/stretchr/testify/require"
)

func BenchmarkMustDeriveKey(b *testing.B) {
	base := make([]byte, 64)
	_, err := rand.Read(base)
	require.NoError(b, err)

	for b.Loop() {
		_ = config.MustDeriveKey(base, "Benchmark", 32)
	}
}
