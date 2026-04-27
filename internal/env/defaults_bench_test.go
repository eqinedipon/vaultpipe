package env

import (
	"fmt"
	"testing"
)

func BenchmarkApplyDefaults_SmallMap(b *testing.B) {
	dst := map[string]string{"A": "1", "B": "2"}
	specs := []DefaultSpec{
		{Key: "B", Value: "overridden"},
		{Key: "C", Value: "3"},
		{Key: "D", Value: "4"},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ApplyDefaults(dst, specs)
	}
}

func BenchmarkApplyDefaults_LargeMap(b *testing.B) {
	dst := make(map[string]string, 50)
	for i := 0; i < 50; i++ {
		dst[fmt.Sprintf("EXISTING_%d", i)] = fmt.Sprintf("val%d", i)
	}
	specs := make([]DefaultSpec, 100)
	for i := 0; i < 100; i++ {
		specs[i] = DefaultSpec{Key: fmt.Sprintf("KEY_%d", i), Value: fmt.Sprintf("default%d", i)}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ApplyDefaults(dst, specs)
	}
}

func BenchmarkDefaultsFromMap(b *testing.B) {
	m := make(map[string]string, 20)
	for i := 0; i < 20; i++ {
		m[fmt.Sprintf("K%d", i)] = fmt.Sprintf("v%d", i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DefaultsFromMap(m)
	}
}
