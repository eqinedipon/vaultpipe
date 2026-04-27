package env

import (
	"testing"
)

func BenchmarkChain_ThreeSteps(b *testing.B) {
	tr := NewTransformer(TrimSpaceTransform)
	chain := NewChain(
		WrapTransformer(tr),
		func(m map[string]string) (map[string]string, error) {
			return SanitizeMap(m), nil
		},
		WrapValidation(RequireKeys("DB_HOST")),
	)

	src := map[string]string{
		"db-host": "  localhost  ",
		"db-port": "  5432  ",
		"api-key": "  secret  ",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = chain.Run(src)
	}
}

func BenchmarkChain_SingleStep(b *testing.B) {
	chain := NewChain(func(m map[string]string) (map[string]string, error) {
		return m, nil
	})
	src := map[string]string{"A": "1", "B": "2", "C": "3"}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = chain.Run(src)
	}
}
