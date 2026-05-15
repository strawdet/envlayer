package aliaser

import "testing"

func BenchmarkApply_SmallMap(b *testing.B) {
	a := New(WithAliases(map[string][]string{
		"DB_HOST": {"DATABASE_HOST"},
		"APP_KEY": {"SERVICE_KEY"},
	}))
	env := map[string]string{
		"DB_HOST": "localhost",
		"APP_KEY": "secret",
		"PORT":    "5432",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = a.Apply(env)
	}
}

func BenchmarkApply_LargeMap(b *testing.B) {
	aliasMap := map[string][]string{}
	env := map[string]string{}
	for i := 0; i < 100; i++ {
		src := fmt.Sprintf("KEY_%d", i)
		tgt := fmt.Sprintf("ALIAS_%d", i)
		aliasMap[src] = []string{tgt}
		env[src] = "value"
	}
	a := New(WithAliases(aliasMap))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = a.Apply(env)
	}
}
