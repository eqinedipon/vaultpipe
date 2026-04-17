package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	vaultAddr  string
	vaultToken string
)

var rootCmd = &cobra.Command{
	Use:   "vaultpipe",
	Short: "Stream secrets from Vault into process environments",
	Long: `vaultpipe fetches secrets from HashiCorp Vault and injects them
into a child process environment without writing secrets to disk.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .vaultpipe.yaml)")
	rootCmd.PersistentFlags().StringVar(&vaultAddr, "vault-addr", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.PersistentFlags().StringVar(&vaultToken, "vault-token", "", "Vault token (overrides VAULT_TOKEN)")
}
