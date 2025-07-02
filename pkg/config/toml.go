package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// LoadFile loads and parses the config file at the given path into the provided struct pointer.
func LoadToml(target interface{}, configfile string) error {
	configfile = filepath.ToSlash(configfile)
	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		return err
	}
	str, err := os.ReadFile(configfile)
	if err != nil {
		return err
	}
	return toml.Unmarshal([]byte(str), target)
}

// SaveConfig saves a Config to the given path.
func SaveToml(path string, cfg interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}
