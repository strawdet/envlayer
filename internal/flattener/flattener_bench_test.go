package flattener_test

import (
	"fmt"
	"testing"

	"github.com/your-org/envlayer/internal/flattener"
)

func BenchmarkFlatten_SmallMap(b *testing.B) {
	f := flattener.New()
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Flatten(env)
	}
}

func BenchmarkFlatten_LargeMap(b *testing.B) {
	f := flattener.New(flattener.WithPrefix("SVC_"))
	env := make(map[string]string, 500)
	for i := 0; i < 500; i++ {
		env[fmt.Sprintf("SVC_KEY_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Flatten(env)
	}
}

func BenchmarkSegments(b *testing.B) {
	f := flattener.New()
	env := make(map[string]string, 200)
	prefixes := []string{"APP", "DB", "SVC", "CACHE", "QUEUE"}
	for i := 0; i < 200; i++ {
		pfx := prefixes[i%len(prefixes)]
		env[fmt.Sprintf("%s_KEY_%d", pfx, i)] = fmt.Sprintf("v%d", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Segments(env)
	}
}
