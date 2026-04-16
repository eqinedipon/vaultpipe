package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockKVv2Server(t *testing.T, mount, secretPath string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/" + mount + "/data/" + secretPath
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"data": payload},
		})
	}))
}

func TestGetSecrets_KVv2Unwrap(t *testing.T) {
	payload := map[string]interface{}{"DB_PASS": "s3cr3t", "API_KEY": "abc123"}
	srv := mockKVv2Server(t, "secret", "myapp", payload)
	defer srv.Close()

	c, err := NewClient(Config{Address: srv.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	sm, err := c.GetSecrets(context.Background(), "secret/myapp")
	if err != nil {
		t.Fatalf("GetSecrets: %v", err)
	}
	if sm["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t, got %q", sm["DB_PASS"])
	}
	if sm["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", sm["API_KEY"])
	}
}

func TestGetMultiple_MergesAndOverrides(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/secret/data/base":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"data": map[string]interface{}{"KEY": "base", "SHARED": "from-base"}},
			})
		case "/v1/secret/data/override":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"data": map[string]interface{}{"SHARED": "from-override"}},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c, _ := NewClient(Config{Address: srv.URL, Token: "tok"})
	sm, err := c.GetMultiple(context.Background(), []string{"secret/base", "secret/override"})
	if err != nil {
		t.Fatalf("GetMultiple: %v", err)
	}
	if sm["KEY"] != "base" {
		t.Errorf("expected KEY=base, got %q", sm["KEY"])
	}
	if sm["SHARED"] != "from-override" {
		t.Errorf("expected SHARED=from-override, got %q", sm["SHARED"])
	}
}

func TestKvPath(t *testing.T) {
	cases := []struct{ in, want string }{
		{"secret/myapp", "secret/data/myapp"},
		{"secret/team/app", "secret/data/team/app"},
		{"secret", "secret"},
	}
	for _, tc := range cases {
		if got := kvPath(tc.in); got != tc.want {
			t.Errorf("kvPath(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
