package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a map of environment variables.
type EnvMap map[string]string

// LoadFile reads a .env file from the given path and returns an EnvMap.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are also ignored.
func LoadFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: open %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("loader: %q line %d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: scan %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a KEY=VALUE line into its components.
func parseLine(line string) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("invalid line %q: missing '='")
	}

	key := strings.TrimSpace(line[:idx])
	if key == "" {
		return "", "", fmt.Errorf("invalid line %q: empty key", line)
	}

	value := strings.TrimSpace(line[idx+1:])
	value = stripQuotes(value)

	return key, value, nil
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
