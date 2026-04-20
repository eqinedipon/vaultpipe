package env_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/env"
)

// TestSnapshotDiffWithInjector verifies that secrets injected via the Injector
// are visible as a diff against a baseline snapshot.
func TestSnapshotDiffWithInjector(t *testing.T) {
	base := map[string]string{
		"HOME": "/home/user",
		"PATH": "/usr/bin",
	}
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "key-xyz",
	}

	before := env.NewSnapshotFromMap(base)

	inj := env.NewInjector(base)
	full := inj.Environ(secrets)

	fullMap := make(map[string]string, len(full))
	for _, e := range full {
		for i, c := range e {
			if c == '=' {
				fullMap[e[:i]] = e[i+1:]
				break
			}
		}
	}

	after := env.NewSnapshotFromMap(fullMap)
	diff := before.Diff(after)

	if diff["DB_PASSWORD"] != "supersecret" {
		t.Errorf("expected DB_PASSWORD in diff, got %q", diff["DB_PASSWORD"])
	}
	if diff["API_KEY"] != "key-xyz" {
		t.Errorf("expected API_KEY in diff, got %q", diff["API_KEY"])
	}
	if _, ok := diff["HOME"]; ok {
		t.Error("HOME should not appear in diff (unchanged)")
	}
}

// TestFilter_RemovesSecretKeysFromSnapshot ensures sensitive keys can be
// stripped before passing the environment to an untrusted subprocess.
func TestFilter_RemovesSecretKeysFromSnapshot(t *testing.T) {
	full := env.NewSnapshotFromMap(map[string]string{
		"APP_HOST":     "localhost",
		"VAULT_TOKEN":  "hvs.supersecret",
		"VAULT_ADDR":   "https://vault:8200",
		"APP_LOG":      "info",
	})

	safe := full.Filter(env.DenyList("VAULT_TOKEN"))

	if _, ok := safe.Get("VAULT_TOKEN"); ok {
		t.Fatal("VAULT_TOKEN must not appear in safe snapshot")
	}
	if v, _ := safe.Get("VAULT_ADDR"); v != "https://vault:8200" {
		t.Errorf("VAULT_ADDR should be retained, got %q", v)
	}
	if v, _ := safe.Get("APP_HOST"); v != "localhost" {
		t.Errorf("APP_HOST should be retained, got %q", v)
	}
}
