package env

import (
	"testing"
)

func TestValidate_AllPresentAndNonEmpty(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	if err := Validate(env, RequireKeys("DB_HOST", "DB_PORT")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
	}
	err := Validate(env, RequireKeys("DB_HOST", "DB_PASSWORD"))
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Missing) != 1 || ve.Missing[0] != "DB_PASSWORD" {
		t.Errorf("unexpected Missing: %v", ve.Missing)
	}
}

func TestValidate_EmptyValueFlagged(t *testing.T) {
	env := map[string]string{
		"API_KEY": "",
		"API_URL": "https://example.com",
	}
	err := Validate(env, RequireKeys("API_KEY"), NoEmptyValues())
	if err == nil {
		t.Fatal("expected error for empty value")
	}
	ve := err.(*ValidationError)
	if len(ve.Invalid) != 1 || ve.Invalid[0] != "API_KEY" {
		t.Errorf("unexpected Invalid: %v", ve.Invalid)
	}
}

func TestValidate_NoEmptyValues_AllKeys(t *testing.T) {
	env := map[string]string{
		"GOOD": "value",
		"BAD":  "",
	}
	err := Validate(env, NoEmptyValues())
	if err == nil {
		t.Fatal("expected error")
	}
	ve := err.(*ValidationError)
	if len(ve.Invalid) != 1 || ve.Invalid[0] != "BAD" {
		t.Errorf("unexpected Invalid: %v", ve.Invalid)
	}
	if len(ve.Missing) != 0 {
		t.Errorf("unexpected Missing: %v", ve.Missing)
	}
}

func TestValidate_ErrorMessage(t *testing.T) {
	env := map[string]string{}
	err := Validate(env, RequireKeys("FOO", "BAR"))
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if msg == "" {
		t.Error("error message must not be empty")
	}
}

func TestIsValidationError(t *testing.T) {
	env := map[string]string{}
	err := Validate(env, RequireKeys("MISSING"))
	if !IsValidationError(err) {
		t.Error("IsValidationError should return true for *ValidationError")
	}
}

func TestValidate_NoOptions_AlwaysPasses(t *testing.T) {
	env := map[string]string{}
	if err := Validate(env); err != nil {
		t.Fatalf("expected nil error with no options, got %v", err)
	}
}
