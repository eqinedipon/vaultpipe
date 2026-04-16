package env

import (
	"os/exec"
	"strings"
	"testing"
)

func TestEnviron_SecretOverridesBase(t *testing.T) {
	base := []string{"PATH=/usr/bin", "HOME=/root", "DB_PASS=old"}
	secrets := map[string]string{"DB_PASS": "new_secret", "API_KEY": "abc123"}

	inj := NewInjector(secrets)
	env := inj.Environ(base)

	envm := toMap(env)

	if envm["DB_PASS"] != "new_secret" {
		t.Errorf("expected DB_PASS=new_secret, got %s", envm["DB_PASS"])
	}
	if envm["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", envm["API_KEY"])
	}
	if envm["PATH"] != "/usr/bin" {
		t.Errorf("expected PATH=/usr/bin, got %s", envm["PATH"])
	}
}

func TestEnviron_EmptyBase(t *testing.T) {
	inj := NewInjector(map[string]string{"SECRET": "val"})
	env := inj.Environ([]string{})
	envm := toMap(env)
	if envm["SECRET"] != "val" {
		t.Errorf("expected SECRET=val, got %s", envm["SECRET"])
	}
}

func TestApplyToCmd(t *testing.T) {
	cmd := exec.Command("env")
	inj := NewInjector(map[string]string{"INJECTED": "yes"})
	inj.ApplyToCmd(cmd, []string{"EXISTING=1"})

	envm := toMap(cmd.Env)
	if envm["INJECTED"] != "yes" {
		t.Error("INJECTED key not found in cmd.Env")
	}
	if envm["EXISTING"] != "1" {
		t.Error("EXISTING key not preserved in cmd.Env")
	}
}

func toMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
