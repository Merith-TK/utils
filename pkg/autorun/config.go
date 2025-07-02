package autorun

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Autorun     string            `toml:"autorun,omitempty"`
	WorkDir     string            `toml:"workDir,omitempty"`
	Isolate     bool              `toml:"isolated,omitempty"`
	Environment map[string]string `toml:"environment,omitempty"`
}

// LoadConfig loads a Config from the given path. If not found, returns an empty config.
func LoadConfig(path string) (Config, error) {
	cfg := Config{Environment: map[string]string{}}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = toml.Unmarshal(data, &cfg)
	return cfg, err
}

// SaveConfig saves a Config to the given path.
func SaveConfig(path string, cfg Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}
