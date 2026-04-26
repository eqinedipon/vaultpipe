package env

import (
	"strings"
	"testing"
)

func BenchmarkTransformer_TrimSpace(b *testing.B) {
	tr := NewTransformer(TrimSpaceTransform())
	src := make(map[string]string, 50)
	for i := 0; i < 50; i++ {
		key := strings.Repeat("K", 8)
		src[key] = "  some secret value  "
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tr.Apply(src)
	}
}

func BenchmarkTransformer_Chain(b *testing.B) {
	upper := func(_, v string) (string, error) { return strings.ToUpper(v), nil }
	tr := NewTransformer(TrimSpaceTransform(), upper)
	src := map[string]string{
		"SECRET_ONE":   "  alpha  ",
		"SECRET_TWO":   "  beta  ",
		"SECRET_THREE": "  gamma  ",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tr.Apply(src)
	}
}
