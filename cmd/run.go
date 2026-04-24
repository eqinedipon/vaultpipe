package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpipe/internal/audit"
	"github.com/your-org/vaultpipe/internal/cache"
	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/env"
	"github.com/your-org/vaultpipe/internal/process"
	"github.com/your-org/vaultpipe/internal/vault"
)

var (
	cfgFile  string
	cacheTTL time.Duration
	strictEnv bool
)

func init() {
	runCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "path to vaultpipe config file")
	runCmd.Flags().DurationVar(&cacheTTL, "cache-ttl", 5*time.Minute, "secret cache TTL (0 disables cache)")
	runCmd.Flags().BoolVar(&strictEnv, "strict", false, "fail if any required env key cannot be resolved")
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
		return fmt.Errorf("load config: %w", err)
	}

	logger := audit.NewLogger(os.Stderr)

	vc, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	cachedClient := vault.NewCachedClient(vc, cache.New(cacheTTL))

	secrets := make(map[string]string)
	for _, s := range cfg.Secrets {
		fetched, err := cachedClient.GetSecrets(cmd.Context(), s.Path, s.Keys)
		if err != nil {
			logger.LogError(s.Path, err)
			return fmt.Errorf("fetch secrets at %s: %w", s.Path, err)
		}
		logger.LogFetch(s.Path, fetched)
		for k, v := range fetched {
			secrets[k] = v
		}
	}

	var resolveOpts []env.ResolveOption
	if strictEnv {
		resolveOpts = append(resolveOpts, env.Strict)
	}

	required := cfg.RequiredKeys()
	base := env.OSMap()
	resolved, err := env.Resolve(required, env.SanitizeMap(secrets), base, resolveOpts...)
	if err != nil {
		return err
	}

	injector := env.NewInjector(base, resolved)
	runner := process.NewRunner(injector)
	logger.LogExec(args[0], args[1:])

	return runner.Run(cmd.Context(), args[0], args[1:]...)
}
