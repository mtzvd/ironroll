// Package env provides a simple .env file loader.
//
// Usage:
//
//	env.Load() // loads .env from current directory
//
// The loader is intentionally simple:
//   - Reads KEY=VALUE pairs, one per line
//   - Ignores empty lines and lines starting with #
//   - Does not override existing environment variables
//   - Silently skips if .env file does not exist
package env

import (
	"bufio"
	"os"
	"strings"
)

// Load reads the .env file and sets environment variables.
// It does not override variables that are already set.
// Returns nil if the file does not exist.
func Load() error {
	return LoadFile(".env")
}

// LoadFile reads the specified file and sets environment variables.
func LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // missing .env is not an error
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first '=' only.
		idx := strings.Index(line, "=")
		if idx < 1 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])

		// Remove surrounding quotes if present.
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// Do not override existing variables.
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
