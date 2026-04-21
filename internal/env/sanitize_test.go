package env

import (
	"testing"
)

func TestSanitizeKey_Uppercase(t *testing.T) {
	if got := SanitizeKey("mykey"); got != "MYKEY" {
		t.Fatalf("expected MYKEY, got %q", got)
	}
}

func TestSanitizeKey_ReplacesDash(t *testing.T) {
	if got := SanitizeKey("my-key"); got != "MY_KEY" {
		t.Fatalf("expected MY_KEY, got %q", got)
	}
}

func TestSanitizeKey_ReplacesSlash(t *testing.T) {
	if got := SanitizeKey("db/password"); got != "DB_PASSWORD" {
		t.Fatalf("expected DB_PASSWORD, got %q", got)
	}
}

func TestSanitizeKey_LeadingDigitReplaced(t *testing.T) {
	got := SanitizeKey("1secret")
	if len(got) == 0 || got[0] == '1' {
		t.Fatalf("leading digit should be replaced, got %q", got)
	}
}

func TestSanitizeKey_Empty(t *testing.T) {
	if got := SanitizeKey(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestSanitizeKey_AlreadyValid(t *testing.T) {
	if got := SanitizeKey("VALID_KEY"); got != "VALID_KEY" {
		t.Fatalf("expected VALID_KEY, got %q", got)
	}
}

func TestSanitizeMap_Keys(t *testing.T) {
	in := map[string]string{
		"my-secret": "val1",
		"db/pass":   "val2",
	}
	out := SanitizeMap(in)
	if v, ok := out["MY_SECRET"]; !ok || v != "val1" {
		t.Errorf("expected MY_SECRET=val1, got %q", v)
	}
	if v, ok := out["DB_PASS"]; !ok || v != "val2" {
		t.Errorf("expected DB_PASS=val2, got %q", v)
	}
}

func TestSanitizeMap_DoesNotMutateInput(t *testing.T) {
	in := map[string]string{"my-key": "v"}
	SanitizeMap(in)
	if _, ok := in["MY_KEY"]; ok {
		t.Fatal("SanitizeMap must not mutate the input map")
	}
}
