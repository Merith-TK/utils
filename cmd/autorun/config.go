package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Merith-TK/utils/pkg/debug"
)

type config struct {
	Autorun     string            `toml:"autorun,omitempty"`
	WorkDir     string            `toml:"workDir,omitempty"`
	Isolate     bool              `toml:"isolated,omitempty"`
	Environment map[string]string `toml:"environment,omitempty"`
}

// TODO: Add custom loader for config replacement keys

func setupConfig(configfile string) (conf *config, err error) {
	debug.Print("setupConfig called with configfile:", configfile)
	// sanitize the config file path
	configfile = filepath.ToSlash(configfile)
	debug.Print("Sanitized configfile path:", configfile)
	// If file does not exist, create it with default values and return the default values
	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		log.Println("Config file", configfile, "does not exist, creating it with default values")
		conf = &config{
			Autorun:     "example.exe",
			WorkDir:     "./",
			Isolate:     false,
			Environment: map[string]string{"FOO": "BAR"},
		}
		debug.Print("Default config created:", conf)
		// Write the config to file
		file, err := os.Create(configfile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		encoder := toml.NewEncoder(file)
		err = encoder.Encode(conf)
		if err != nil {
			return nil, err
		}
		debug.Print("Default config written to file:", configfile)
	}

	// Read and unmarshal the config file
	str, err := os.ReadFile(configfile)
	if err != nil {
		return nil, err
	}
	debug.Print("Config file read:", string(str))
	err = toml.Unmarshal([]byte(str), &conf)
	if err != nil {
		return nil, err
	}
	debug.Print("Config unmarshaled:", conf)
	return conf, nil
}

// Setup the Environment part of the config,
func setupEnvironment(conf *config) (config *config) {
	debug.Print("setupEnvironment called with config:", conf)
	// Variables for env replacement
	drivePath, _ := filepath.Abs("/")
	drivePath = filepath.ToSlash(drivePath)
	drivePath = strings.TrimSuffix(drivePath, "/")
	configEnvReplace := map[string]string{
		"{work}":  conf.WorkDir,
		"{drive}": drivePath,
	}
	debug.Print("Config environment replacements:", configEnvReplace)

	// TODO: Add custom loader for config replacement keys

	// Replace Normal Config options
	for key, value := range configEnvReplace {
		if strings.Contains(conf.Autorun, key) {
			conf.Autorun = filepath.ToSlash(strings.ReplaceAll(conf.Autorun, key, value))
			debug.Print("Replaced Autorun key:", key, "with value:", value, "result:", conf.Autorun)
		}
		if strings.Contains(conf.WorkDir, key) {
			conf.WorkDir = filepath.ToSlash(strings.ReplaceAll(conf.WorkDir, key, value))
			debug.Print("Replaced WorkDir key:", key, "with value:", value, "result:", conf.WorkDir)
		}
	}

	// Replace Environment Variables
	for k, v := range conf.Environment {
		for key, value := range configEnvReplace {
			if strings.Contains(v, key) {
				v = strings.ReplaceAll(v, key, value)
				v = filepath.ToSlash(v)
				debug.Print("Replaced environment variable:", k, "key:", key, "with value:", value, "result:", v)
			}
		}
		os.Setenv(k, v)
	}

	// TODO: Use environment variables to replace config values with {$ENV_VAR}

	return conf
}
