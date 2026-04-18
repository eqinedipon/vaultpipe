package template_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpipe/internal/template"
)

// TestRenderAll_MultipleSecrets ensures several secrets can be composed
// across multiple env-var values in a single RenderAll call.
func TestRenderAll_MultipleSecrets(t *testing.T) {
	secrets := map[string]string{
		"PG_USER": "admin",
		"PG_PASS": "hunter2",
		"PG_HOST": "db.internal",
	}
	r := template.New(secrets)
	pairs := map[string]string{
		"DATABASE_URL": "postgres://{{ .PG_USER }}:{{ .PG_PASS }}@{{ .PG_HOST }}/app",
		"APP_ENV":      "production",
	}
	out, err := r.RenderAll(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://admin:hunter2@db.internal/app"
	if out["DATABASE_URL"] != want {
		t.Errorf("DATABASE_URL: got %q, want %q", out["DATABASE_URL"], want)
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q", out["APP_ENV"])
	}
}

// TestRender_SecretsNotMutated ensures the renderer does not expose the
// internal secrets map to callers via the rendered output unexpectedly.
func TestRender_SecretsNotMutated(t *testing.T) {
	original := map[string]string{"KEY": "value"}
	r := template.New(original)
	original["KEY"] = "mutated"
	out, err := r.Render("{{ .KEY }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "mutated") {
		t.Errorf("renderer reflected external mutation: got %q", out)
	}
}
