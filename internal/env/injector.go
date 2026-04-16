package env

import (
	"fmt"
	"os/exec"
	"strings"
)

// Injector holds secrets to be injected into a subprocess environment.
type Injector struct {
	secrets map[string]string
}

// NewInjector creates an Injector from a map of secret key/value pairs.
func NewInjector(secrets map[string]string) *Injector {
	return &Injector{secrets: secrets}
}

// Environ returns the merged environment: current OS env overridden by secrets.
func (inj *Injector) Environ(base []string) []string {
	merged := make(map[string]string, len(base))
	for _, e := range base {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			merged[parts[0]] = parts[1]
		}
	}
	for k, v := range inj.secrets {
		merged[k] = v
	}
	result := make([]string, 0, len(merged))
	for k, v := range merged {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}

// ApplyToCmd sets the Env field on an exec.Cmd, merging secrets into the
// provided base environment slice (typically os.Environ()).
func (inj *Injector) ApplyToCmd(cmd *exec.Cmd, base []string) {
	cmd.Env = inj.Environ(base)
}
