package env

import "fmt"

// DefaultSpec describes a single default value entry.
type DefaultSpec struct {
	Key      string
	Value    string
	Override bool // if true, overwrite existing value
}

// DefaultsOption configures ApplyDefaults behaviour.
type DefaultsOption func(*defaultsConfig)

type defaultsConfig struct {
	override bool
}

// OverrideExisting makes ApplyDefaults overwrite keys that already have a value.
func OverrideExisting() DefaultsOption {
	return func(c *defaultsConfig) { c.override = true }
}

// ApplyDefaults merges a slice of DefaultSpec entries into dst.
// By default it does NOT overwrite keys that are already present in dst.
// Pass OverrideExisting() to change that behaviour.
// It returns an error if any spec has an empty Key.
func ApplyDefaults(dst map[string]string, specs []DefaultSpec, opts ...DefaultsOption) (map[string]string, error) {
	cfg := &defaultsConfig{}
	for _, o := range opts {
		o(cfg)
	}

	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	for _, s := range specs {
		if s.Key == "" {
			return nil, fmt.Errorf("env/defaults: empty key in DefaultSpec")
		}
		_, exists := out[s.Key]
		if !exists || cfg.override {
			out[s.Key] = s.Value
		}
	}
	return out, nil
}

// DefaultsFromMap converts a plain map into a []DefaultSpec slice.
// All entries are non-override by default.
func DefaultsFromMap(m map[string]string) []DefaultSpec {
	specs := make([]DefaultSpec, 0, len(m))
	for k, v := range m {
		specs = append(specs, DefaultSpec{Key: k, Value: v})
	}
	return specs
}
