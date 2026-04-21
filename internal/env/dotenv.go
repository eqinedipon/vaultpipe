// Package env provides utilities for environment variable management.
package env

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ParseDotEnv reads key=value pairs from r, ignoring blank lines and
// lines that begin with '#'. Inline comments are not supported.
// Quoted values (single or double) have their surrounding quotes stripped.
func ParseDotEnv(r io.Reader) (map[string]string, error) {
	out := make(map[string]string)
	scanner := bufio.NewScanner(r)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			return nil, fmt.Errorf("dotenv: invalid syntax on line %d: %q", lineNo, line)
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = stripQuotes(val)
		out[key] = val
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dotenv: scanner error: %w", err)
	}
	return out, nil
}

// LoadDotEnvFile opens path and delegates to ParseDotEnv.
// It returns an empty map (and no error) when the file does not exist so
// callers can treat a missing .env file as optional.
func LoadDotEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dotenv: open %s: %w", path, err)
	}
	defer f.Close()
	return ParseDotEnv(f)
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
