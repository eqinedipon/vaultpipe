// Package env provides utilities for injecting Vault secrets into subprocess
// environments. It merges secret key/value pairs with a base environment
// (typically os.Environ()), ensuring secrets override any existing values
// without persisting them to disk or the parent process environment.
//
// Usage:
//
//	inj := env.NewInjector(secrets)
//	inj.ApplyToCmd(cmd, os.Environ())
//	err := cmd.Run()
package env
