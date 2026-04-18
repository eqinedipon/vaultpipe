package template_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/template"
)

func TestRender_SimpleSubstitution(t *testing.T) {
	r := template.New(map[string]string{"DB_PASS": "s3cr3t"})
	out, err := r.Render("password={{ .DB_PASS }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "password=s3cr3t" {
		t.Errorf("got %q", out)
	}
}

func TestRender_MissingKey_ReturnsError(t *testing.T) {
	r := template.New(map[string]string{})
	_, err := r.Render("val={{ .MISSING }}")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRender_NoTemplate_PassThrough(t *testing.T) {
	r := template.New(map[string]string{})
	out, err := r.Render("plain string")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "plain string" {
		t.Errorf("got %q", out)
	}
}

func TestRenderAll_InterpolatesOnlyTemplated(t *testing.T) {
	r := template.New(map[string]string{"TOKEN": "abc123"})
	input := map[string]string{
		"AUTH": "Bearer {{ .TOKEN }}",
		"PLAIN": "no-template",
	}
	out, err := r.RenderAll(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["AUTH"] != "Bearer abc123" {
		t.Errorf("AUTH: got %q", out["AUTH"])
	}
	if out["PLAIN"] != "no-template" {
		t.Errorf("PLAIN: got %q", out["PLAIN"])
	}
}

func TestRenderAll_MissingKey_ReturnsError(t *testing.T) {
	r := template.New(map[string]string{})
	_, err := r.RenderAll(map[string]string{"X": "{{ .NOPE }}"})
	if err == nil {
		t.Fatal("expected error")
	}
}
