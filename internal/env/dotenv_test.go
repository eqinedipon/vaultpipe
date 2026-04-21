package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseDotEnv_BasicPairs(t *testing.T) {
	input := "FOO=bar\nBAZ=qux\n"
	m, err := ParseDotEnv(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestParseDotEnv_IgnoresComments(t *testing.T) {
	input := "# comment\nKEY=value\n"
	m, err := ParseDotEnv(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 1 || m["KEY"] != "value" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestParseDotEnv_StripsDoubleQuotes(t *testing.T) {
	m, err := ParseDotEnv(strings.NewReader(`SECRET="hello world"`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["SECRET"] != "hello world" {
		t.Fatalf("got %q", m["SECRET"])
	}
}

func TestParseDotEnv_StripsSingleQuotes(t *testing.T) {
	m, err := ParseDotEnv(strings.NewReader("TOKEN='abc123'"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["TOKEN"] != "abc123" {
		t.Fatalf("got %q", m["TOKEN"])
	}
}

func TestParseDotEnv_InvalidLine(t *testing.T) {
	_, err := ParseDotEnv(strings.NewReader("NODEQUALS\n"))
	if err == nil {
		t.Fatal("expected error for line without '='")
	}
}

func TestParseDotEnv_EmptyValue(t *testing.T) {
	m, err := ParseDotEnv(strings.NewReader("EMPTY=\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := m["EMPTY"]; !ok || v != "" {
		t.Fatalf("expected empty string, got %q", v)
	}
}

func TestLoadDotEnvFile_MissingFile_ReturnsEmpty(t *testing.T) {
	m, err := LoadDotEnvFile("/nonexistent/.env")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(m) != 0 {
		t.Fatalf("expected empty map, got %v", m)
	}
}

func TestLoadDotEnvFile_ReadsFromDisk(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("DISK_KEY=disk_val\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := LoadDotEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["DISK_KEY"] != "disk_val" {
		t.Fatalf("unexpected map: %v", m)
	}
}
