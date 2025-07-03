# config

This package provides configuration loading and environment setup utilities.

## Functions

### EnvKeyReplace

```
func EnvKeyReplace(input string, replacements map[string]string) string
```
Replaces all {key} in the input string with their values from the replacements map.

### EnvOverride

```
func EnvOverride(env map[string]string)
```
Sets environment variables from the provided map.

### LoadToml

```
func LoadToml(target interface{}, configfile string) error
```
Loads TOML configuration from a file into the provided struct pointer.

### SaveToml

```
func SaveToml(path string, cfg interface{}) error
```
Saves a struct as TOML to the given file path.

## Example

```go
import "github.com/Merith-TK/utils/pkg/config"

type MyConfig struct {
    Name string `toml:"name"`
    Port int    `toml:"port"`
}

var cfg MyConfig
err := config.LoadToml(&cfg, "config.toml")
if err != nil {
    // handle error
}

cfg.Name = "test"
err = config.SaveToml("config.toml", cfg)
if err != nil {
    // handle error
}
``` 