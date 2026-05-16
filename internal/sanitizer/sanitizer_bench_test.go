package sanitizer_test

import (
	"fmt"
	"testing"

	"github.com/your-org/envlayer/internal/sanitizer"
)

func BenchmarkApply_SmallMap(b *testing.B) {
	s := sanitizer.New(sanitizer.WithUpperKeys(), sanitizer.WithStripControlChars())
	env := map[string]string{
		"app-name":  "envlayer",
		"db host":   "localhost",
		"PORT":      "5432",
		"LOG_LEVEL": "info",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Apply(env)
	}
}

func BenchmarkApply_LargeMap(b *testing.B) {
	s := sanitizer.New(sanitizer.WithUpperKeys(), sanitizer.WithStripControlChars())
	env := make(map[string]string, 200)
	for i := 0; i < 200; i++ {
		env[fmt.Sprintf("key-%d", i)] = fmt.Sprintf("value\x01%d", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Apply(env)
	}
}
