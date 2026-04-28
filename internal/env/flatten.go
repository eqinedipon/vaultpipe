// Package env provides utilities for environment variable manipulation.
package env

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested maps are flattened into env-style keys.
type FlattenOptions struct {
	// Separator is placed between key segments. Defaults to "_".
	Separator string
	// Prefix is prepended to every resulting key.
	Prefix string
	// UpperCase converts all keys to upper case after flattening.
	UpperCase bool
}

var defaultFlattenOptions = FlattenOptions{
	Separator: "_",
	UpperCase: true,
}

// FlattenMap recursively flattens a nested map[string]any into a flat
// map[string]string suitable for use as environment variables.
//
//	FlattenMap(map[string]any{"db": map[string]any{"host": "localhost", "port": 5432}}, FlattenOptions{})
//	// => {"DB_HOST": "localhost", "DB_PORT": "5432"}
func FlattenMap(src map[string]any, opts FlattenOptions) (map[string]string, error) {
	if opts.Separator == "" {
		opts.Separator = defaultFlattenOptions.Separator
	}
	out := make(map[string]string)
	if err := flattenInto(src, opts.Prefix, opts.Separator, opts.UpperCase, out); err != nil {
		return nil, err
	}
	return out, nil
}

// FlattenKeys returns a sorted slice of all keys produced by FlattenMap.
func FlattenKeys(src map[string]any, opts FlattenOptions) ([]string, error) {
	m, err := FlattenMap(src, opts)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

func flattenInto(src map[string]any, prefix, sep string, upper bool, dst map[string]string) error {
	for k, v := range src {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + sep + k
		}
		if upper {
			fullKey = strings.ToUpper(fullKey)
		}
		switch val := v.(type) {
		case map[string]any:
			if err := flattenInto(val, fullKey, sep, upper, dst); err != nil {
				return err
			}
		case nil:
			dst[fullKey] = ""
		default:
			dst[fullKey] = fmt.Sprintf("%v", val)
		}
	}
	return nil
}
