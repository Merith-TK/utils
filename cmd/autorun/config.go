package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type config struct {
	Autorun     string            `toml:"autorun,omitempty"`
	WorkDir     string            `toml:"workDir,omitempty"`
	Isolate     bool              `toml:"isolated,omitempty"`
	Environment map[string]string `toml:"environment,omitempty"`
}

// TODO: Add custom loader for config replacement keys

func setupConfig(configfile string) (conf *config, err error) {
	// sanitize the config file path
	configfile = filepath.ToSlash(configfile)
	// If file does not exist, create it with default values and return the default values
	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		log.Println("Config file", configfile, "does not exist, creating it with default values")
		conf = &config{
			Autorun:     "example.exe",
			WorkDir:     "./",
			Isolate:     false,
			Environment: map[string]string{"FOO": "BAR"},
		}
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
	}

	// Read and unmarshal the config file
	str, err := os.ReadFile(configfile)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal([]byte(str), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// Setup the Environment part of the config,
func setupEnvironment(conf *config) (config *config) {
	// Variables for env replacement
	drivePath, _ := filepath.Abs("/")
	drivePath = filepath.ToSlash(drivePath)
	drivePath = strings.TrimSuffix(drivePath, "/")
	configEnvReplace := map[string]string{
		"{work}":  conf.WorkDir,
		"{drive}": drivePath,
	}

	// TODO: Add custom loader for config replacement keys

	// Replace Normal Config options
	for key, value := range configEnvReplace {
		if strings.Contains(conf.Autorun, key) {
			conf.Autorun = filepath.ToSlash(strings.ReplaceAll(conf.Autorun, key, value))
		}
		if strings.Contains(conf.WorkDir, key) {
			conf.WorkDir = filepath.ToSlash(strings.ReplaceAll(conf.WorkDir, key, value))
		}
	}

	// Replace Environment Variables
	for k, v := range conf.Environment {
		for key, value := range configEnvReplace {
			if strings.Contains(v, key) {
				v = strings.ReplaceAll(v, key, value)
				v = filepath.ToSlash(v)
			}
		}
		os.Setenv(k, v)
	}

	// TODO: Use environment variables to replace config values with {$ENV_VAR}

	return conf
}
