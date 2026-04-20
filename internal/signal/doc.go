// Package signal implements OS signal forwarding for vaultpipe.
//
// When vaultpipe wraps a subprocess it must relay signals (e.g. SIGINT from
// Ctrl-C, SIGTERM from a container orchestrator, SIGHUP for config reload)
// so that the child process can perform its own graceful shutdown rather than
// being abandoned while vaultpipe exits.
//
// Usage:
//
//	f := signal.New(cmd.Process)
//	f.Start()
//	defer f.Stop()
//
// By default SIGINT, SIGTERM, and SIGHUP are forwarded. Pass explicit signals
// to New to override this set.
package signal
