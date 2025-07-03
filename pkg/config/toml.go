package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// LoadToml loads and parses the config file at the given path into the provided struct pointer.
// The config file must be in TOML format.
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

// SaveToml saves a struct as TOML to the given file path.
func SaveToml(path string, cfg interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}
