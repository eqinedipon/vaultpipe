package cmd

import (
	"bytes"
	"testing"
)

func TestRunCmd_MissingArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"run"})
	var buf bytes.Buffer
	rootCmd.SetErr(&buf)

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no command provided")
	}
}

func TestRunCmd_Registered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "run -- <command> [args...]" {
			return
		}
	}
	t.Fatal("run subcommand not registered")
}

func TestRootCmd_PersistentFlags(t *testing.T) {
	flags := []string{"config", "vault-addr", "vault-token"}
	for _, f := range flags {
		if rootCmd.PersistentFlags().Lookup(f) == nil {
			t.Errorf("expected persistent flag %q to be registered", f)
		}
	}
}
