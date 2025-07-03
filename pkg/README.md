# Package Documentation

This directory contains shared Go packages used across various utilities in this repository.

## Packages

### archive
Package archive provides utilities for working with archive formats.

**Functions:**
- `Unzip(src, dest string) error` - Extracts a ZIP archive from src to dest directory

### autorun  
Package autorun provides configuration structures and utilities for autorun functionality.

**Types:**
- `Config` - Configuration structure for autorun settings
  - `Autorun string` - Path to executable to run
  - `WorkDir string` - Working directory for execution
  - `Isolate bool` - Whether to run in isolated environment
  - `Environment map[string]string` - Custom environment variables

**Functions:**
- `LoadConfig(path string) (Config, error)` - Loads configuration from TOML file
- `SaveConfig(path string, cfg Config) error` - Saves configuration to TOML file

### config
Package config provides configuration loading and environment setup for autorun and other utilities.

**Functions:**
- `EnvKeyReplace(input string, replacements map[string]string) string` - Replaces placeholders in strings with values
- `EnvOverride(env map[string]string)` - Sets environment variables from map
- `LoadToml(target interface{}, configfile string) error` - Loads TOML configuration into struct
- `SaveToml(path string, cfg interface{}) error` - Saves struct to TOML file

### debug
Package debug provides utilities for debugging purposes with conditional output and stacktraces.

**Global Variables:**
- `Title string` - Current debug title prefix

**Functions:**
- `Print(message ...any)` - Prints debug message if debug mode enabled
- `Println(message ...any)` - Prints debug message with newline
- `SetTitle(title string)` - Sets debug message title prefix
- `GetTitle() string` - Gets current debug title
- `ResetTitle()` - Resets title to default
- `SetDebug(enabled bool)` - Toggles debug mode
- `GetDebug() bool` - Gets debug mode status
- `SetStacktrace(enabled bool)` - Toggles stacktrace mode
- `GetStacktrace() bool` - Gets stacktrace mode status
- `Suicide(timeout int)` - Self-destructs after timeout if suicide mode enabled

**Flags:**
- `-debug` - Enable debug mode (or set DEBUG=true env var)
- `-stacktrace` - Enable stacktraces (or set STACKTRACE=true env var)  
- `-suicide` - Enable suicide mode (or set SUICIDE=true env var)

### driveutil
Package driveutil provides drive enumeration, metadata extraction, and utility functions for Windows drives.

**Types:**
- `DriveStore map[string]bool` - Tracks detected drives by unique ID
- `DriveInfo` - Information about a drive
  - `Letter string` - Drive letter (e.g., "C:\\")
  - `Label string` - Volume label
  - `Serial uint32` - Volume serial number
  - `Type uint32` - Drive type (DRIVE_FIXED, DRIVE_REMOVABLE, etc.)

**Functions:**
- `(store DriveStore) DetectDrives(onNewDrive func(drive string, serial uint32))` - Detects new drives and calls callback
- `(store DriveStore) MonitorDrives(onNewDrive func(drive string, serial uint32), interval time.Duration)` - Continuously monitors for new drives
- `GetVolumeSerialNumber(root string) (uint32, error)` - Gets volume serial number for drive
- `DriveExists(drive string) bool` - Checks if drive path exists
- `ListDrives() []DriveInfo` - Returns slice of all available drives

## Usage

Import packages as needed:

```go
import (
    "github.com/Merith-TK/utils/pkg/config"
    "github.com/Merith-TK/utils/pkg/debug"
    "github.com/Merith-TK/utils/pkg/driveutil"
)
```

## Examples

### Debug Package
```go
// Enable debug mode
debug.SetDebug(true)
debug.SetTitle("MyApp")

// Print debug messages
debug.Print("Starting application")
debug.Println("Configuration loaded")

// Reset title when done
defer debug.ResetTitle()
```

### Config Package
```go
type MyConfig struct {
    Name string `toml:"name"`
    Port int    `toml:"port"`
}

var cfg MyConfig
err := config.LoadToml(&cfg, "config.toml")
if err != nil {
    log.Fatal(err)
}
```

### DriveUtil Package
```go
// List all drives
drives := driveutil.ListDrives()
for _, drive := range drives {
    fmt.Printf("Drive: %s, Label: %s, Serial: %08X\n", 
        drive.Letter, drive.Label, drive.Serial)
}

// Monitor for new drives
store := driveutil.DriveStore{}
store.MonitorDrives(func(drive string, serial uint32) {
    fmt.Printf("New drive detected: %s (Serial: %08X)\n", drive, serial)
}, 5*time.Second)
```
