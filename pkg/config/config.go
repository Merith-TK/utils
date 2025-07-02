// Package config provides configuration loading and environment setup for autorun and other utilities.
package config

import (
	"os"
	"strings"
)

// EnvKeyReplace replaces all {key} in the input string with their values from the replacements map.
func EnvKeyReplace(input string, replacements map[string]string) string {
	for key, value := range replacements {
		input = strings.ReplaceAll(input, key, value)
	}
	return input
}

// EnvOverride sets environment variables from the provided map.
func EnvOverride(env map[string]string) {
	for k, v := range env {
		os.Setenv(k, v)
	}
}
