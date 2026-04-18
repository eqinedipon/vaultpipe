// Package template provides Go-template-based interpolation of Vault secret
// values into arbitrary string pairs.
//
// Typical usage: render environment variable values that reference secrets
// fetched from Vault, using standard Go template syntax:
//
//	{{ .SECRET_KEY }}
//
// The Renderer treats missing keys as errors to prevent silent
// misconfigurations from reaching child processes.
package template
