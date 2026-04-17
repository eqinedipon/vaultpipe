package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpipe/vaultpipe/internal/config"
	"github.com/vaultpipe/vaultpipe/internal/env"
	"github.com/vaultpipe/vaultpipe/internal/process"
	"github.com/vaultpipe/vaultpipe/internal/vault"
)

var runCmd = &cobra.Command{
	Use:   "run -- <command> [args...]",
	Short: "Run a command with secrets injected into its environment",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runCommand,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if vaultAddr != "" {
		cfg.Vault.Address = vaultAddr
	}
	if vaultToken != "" {
		cfg.Vault.Token = vaultToken
	}

	vc, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	secrets, err := vc.GetSecrets(cmd.Context(), cfg.Secrets)
	if err != nil {
		return fmt.Errorf("fetching secrets: %w", err)
	}

	injector := env.NewInjector(os.Environ(), secrets)
	runner := process.NewRunner(args[0], args[1:], injector.Environ())

	return runner.Run(cmd.Context())
}
