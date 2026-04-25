package env

import "testing"

var sink map[string]string

func BenchmarkCoerceMap_Mixed(b *testing.B) {
	input := map[string]any{
		"KEY_STR":   "some-secret-value",
		"KEY_INT":   int(12345),
		"KEY_FLOAT": float64(9.99),
		"KEY_BOOL":  true,
		"KEY_NIL":   nil,
		"KEY_I64":   int64(1<<40),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := CoerceMap(input)
		if err != nil {
			b.Fatal(err)
		}
		sink = out
	}
}

func BenchmarkCoerceValue_String(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = CoerceValue("plain-string-value")
	}
}

func BenchmarkCoerceValue_Float64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = CoerceValue(float64(3.141592653589793))
	}
}
