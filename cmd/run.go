package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/env"
	"github.com/your-org/vaultpipe/internal/process"
	"github.com/your-org/vaultpipe/internal/vault"
)

var (
	cfgFile    string
	cacheTTL   time.Duration
	sanitize   bool
)

func init() {
	runCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "path to vaultpipe config file")
	runCmd.Flags().DurationVar(&cacheTTL, "cache-ttl", 5*time.Minute, "secret cache TTL (0 disables cache)")
	runCmd.Flags().BoolVar(&sanitize, "sanitize-keys", true, "sanitize secret keys to valid POSIX env var names")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run -- <command> [args...]",
	Short: "Run a command with secrets injected into its environment",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runCommand,
}

func runCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(vault.ClientConfig{
		Address: cfg.VaultAddr,
		Token:   cfg.VaultToken,
	})
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	secrets, err := client.GetSecrets(cmd.Context(), cfg.Secrets)
	if err != nil {
		return fmt.Errorf("fetching secrets: %w", err)
	}

	if sanitize {
		secrets = env.SanitizeMap(secrets)
	}

	injector := env.NewInjector(os.Environ(), secrets)
	runner := process.NewRunner()

	return runner.Run(cmd.Context(), args, injector.Environ())
}
