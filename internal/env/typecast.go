package env

import (
	"fmt"
	"strconv"
	"strings"
)

// TypeCastRule defines how a specific environment variable key should be
// interpreted when casting to a target type string ("bool", "int", "float").
type TypeCastRule struct {
	Key      string
	TypeName string // "bool", "int", "float", "string"
}

// TypeCaster applies type-aware casting rules to a string-valued env map,
// returning a new map with values normalized to their canonical string form.
// For example, "true", "1", "yes" all normalize to "true" under a bool rule.
type TypeCaster struct {
	rules map[string]string // key -> type name
}

// NewTypeCaster constructs a TypeCaster from the provided rules.
func NewTypeCaster(rules []TypeCastRule) *TypeCaster {
	m := make(map[string]string, len(rules))
	for _, r := range rules {
		m[strings.ToUpper(r.Key)] = strings.ToLower(r.TypeName)
	}
	return &TypeCaster{rules: m}
}

// Cast applies casting rules to src, returning a new normalized map.
// Keys without a matching rule are passed through unchanged.
// Returns an error if a value cannot be cast to its declared type.
func (tc *TypeCaster) Cast(src map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		typeName, ok := tc.rules[strings.ToUpper(k)]
		if !ok {
			out[k] = v
			continue
		}
		normalized, err := castValue(v, typeName)
		if err != nil {
			return nil, fmt.Errorf("typecast: key %q value %q: %w", k, v, err)
		}
		out[k] = normalized
	}
	return out, nil
}

// castValue normalizes v to the canonical string form for typeName.
func castValue(v, typeName string) (string, error) {
	switch typeName {
	case "bool":
		b, err := parseBool(v)
		if err != nil {
			return "", err
		}
		if b {
			return "true", nil
		}
		return "false", nil
	case "int":
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as int", v)
		}
		return strconv.FormatInt(int64(f), 10), nil
	case "float":
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as float", v)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case "string":
		return v, nil
	default:
		return "", fmt.Errorf("unknown type %q", typeName)
	}
}

// parseBool accepts truthy/falsy strings beyond strconv.ParseBool.
func parseBool(v string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1", "yes", "on", "enabled":
		return true, nil
	case "false", "0", "no", "off", "disabled", "":
		return false, nil
	}
	return false, fmt.Errorf("cannot parse %q as bool", v)
}
