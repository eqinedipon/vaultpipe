package env_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/env"
)

// TestSanitizeMap_ThenInject verifies that sanitized secret keys can be
// injected into a process environment without collision with base env vars.
func TestSanitizeMap_ThenInject(t *testing.T) {
	base := []string{"PATH=/usr/bin", "HOME=/root"}

	raw := map[string]string{
		"db/password":    "s3cr3t",
		"api-key":        "tok123",
		"already_UPPER":  "upper",
	}

	sanitized := env.SanitizeMap(raw)

	injector := env.NewInjector(base, sanitized)
	environ := injector.Environ()

	m := toEnvMap(environ)

	if m["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("DB_PASSWORD not injected correctly: %q", m["DB_PASSWORD"])
	}
	if m["API_KEY"] != "tok123" {
		t.Errorf("API_KEY not injected correctly: %q", m["API_KEY"])
	}
	if m["ALREADY_UPPER"] != "upper" {
		t.Errorf("ALREADY_UPPER not injected correctly: %q", m["ALREADY_UPPER"])
	}
	// Base vars must still be present.
	if m["PATH"] != "/usr/bin" {
		t.Errorf("PATH should be preserved, got %q", m["PATH"])
	}
}

func toEnvMap(environ []string) map[string]string {
	m := make(map[string]string, len(environ))
	for _, e := range environ {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				m[e[:i]] = e[i+1:]
				break
			}
		}
	}
	return m
}
