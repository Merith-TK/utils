// Package config provides configuration management utilities including environment variable
// manipulation, key replacement, and TOML configuration file handling.
//
// The package offers two main areas of functionality:
//   - Environment variable management with key-value replacement
//   - TOML configuration file loading and saving
//
// Key replacement allows templating in configuration strings using {KEY} syntax.
// Environment override provides a convenient way to set multiple environment variables.
//
// Example usage:
//
//	// Key replacement
//	replacements := map[string]string{"{USER}": "john", "{HOME}": "/home/john"}
//	result := config.EnvKeyReplace("User: {USER}, Home: {HOME}", replacements)
//
//	// Environment override
//	env := map[string]string{"DEBUG": "true", "PORT": "8080"}
//	config.EnvOverride(env)
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
