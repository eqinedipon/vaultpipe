// Package template provides secret interpolation into string templates.
package template

import (
	"bytes"
	"fmt"
	"strings"
	gotemplate "text/template"
)

// Renderer interpolates secret values into Go templates.
type Renderer struct {
	secrets map[string]string
}

// New creates a Renderer with the provided secret map.
func New(secrets map[string]string) *Renderer {
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	return &Renderer{secrets: copy}
}

// Render interpolates the template string using the secret map.
// Template variables use Go template syntax: {{ .SECRET_KEY }}
func (r *Renderer) Render(tmpl string) (string, error) {
	t, err := gotemplate.New("").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("template parse: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, r.secretsAsMap()); err != nil {
		return "", fmt.Errorf("template execute: %w", err)
	}
	return buf.String(), nil
}

// RenderAll renders all values in the provided string map and returns
// a new map with interpolated values.
func (r *Renderer) RenderAll(pairs map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(pairs))
	for k, v := range pairs {
		if !strings.Contains(v, "{{") {
			out[k] = v
			continue
		}
		rendered, err := r.Render(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = rendered
	}
	return out, nil
}

func (r *Renderer) secretsAsMap() map[string]string {
	return r.secrets
}
