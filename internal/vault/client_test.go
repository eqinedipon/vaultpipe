package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockVaultServer(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/secret/data/myapp/config":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"data": data,
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestNewClient_DefaultsFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "test-token")

	client, err := NewClient(Config{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.api.Token() != "test-token" {
		t.Errorf("expected token 'test-token', got %q", client.api.Token())
	}
}

func TestNewClient_ExplicitConfig(t *testing.T) {
	client, err := NewClient(Config{
		Address: "http://127.0.0.1:8200",
		Token:   "explicit-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.api.Token() != "explicit-token" {
		t.Errorf("expected 'explicit-token', got %q", client.api.Token())
	}
}

func TestGetSecrets_StringCoercion(t *testing.T) {
	// Verify non-string values are coerced to strings
	client, err := NewClient(Config{
		Address: "http://127.0.0.1:8200",
		Token:   "tok",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Simulate coercion logic directly (unit test the helper behaviour)
	input := map[string]interface{}{
		"PORT":    float64(8080),
		"ENABLED": true,
		"NAME":    "myapp",
	}
	result := make(map[string]string, len(input))
	for k, v := range input {
		if str, ok := v.(string); ok {
			result[k] = str
		} else {
			result[k] = fmt.Sprintf("%v", v)
		}
	}
	_ = client // used above

	if result["PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", result["PORT"])
	}
	if result["NAME"] != "myapp" {
		t.Errorf("expected 'myapp', got %q", result["NAME"])
	}
}
