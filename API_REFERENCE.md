# API Reference

Complete API documentation for all packages in the utils library.

## Package: debug

**Import**: `github.com/Merith-TK/utils/pkg/debug`

### Overview

The debug package provides comprehensive debugging utilities with conditional logging, stacktrace output, and self-destruct functionality for development and testing.

### Configuration

The package can be configured via command-line flags or environment variables:

| Flag | Environment | Description |
|------|-------------|-------------|
| `-debug` | `DEBUG=true` | Enable debug output |
| `-stacktrace` | `STACKTRACE=true` | Enable stacktrace in debug output |
| `-suicide` | `SUICIDE=true` | Enable self-destruct functionality |

### Functions

#### SetDebug
```go
func SetDebug(enabled bool)
```
Programmatically enables or disables debug mode, overriding any command-line flag or environment variable settings.

#### GetDebug
```go
func GetDebug() bool
```
Returns true if debug mode is currently enabled, either through flags, environment variables, or SetDebug calls.

#### SetTitle
```go
func SetTitle(title string)
```
Sets a custom title prefix for debug messages globally. This affects all subsequent debug output until reset or changed. The title appears in debug output as `[DEBUG] {title} message`. Best practice is to use `defer ResetTitle()` after calling this function.

#### GetTitle
```go
func GetTitle() string
```
Returns the currently set debug message title prefix. Returns an empty string if no custom title has been set.

#### ResetTitle
```go
func ResetTitle()
```
Resets the debug message title prefix to the default empty value. This should be called to clean up after using SetTitle, typically with defer.

#### Print
```go
func Print(message ...any)
```
Outputs the given message to standard output if debug mode is enabled. Messages are prefixed with `[DEBUG]` and optionally include a custom title if set. If stacktrace mode is also enabled, a filtered stack trace is included that excludes internal debug package frames for cleaner output.

#### Println
```go
func Println(message ...any)
```
Prints the given message to the standard output followed by a newline character, but only if debug mode is enabled. This is a convenience wrapper around Print that automatically adds a newline.

#### Suicide
```go
func Suicide(timeout int)
```
Enables a self-destruct mechanism that will terminate the process after the specified timeout in seconds, but only if suicide mode is enabled via flag or environment variable. This is primarily used for testing and development to prevent runaway processes. The function is non-blocking and starts a goroutine to handle the timeout.

#### SetStacktrace
```go
func SetStacktrace(enabled bool)
```
Programmatically enables or disables stacktrace output in debug messages, overriding any command-line flag or environment variable settings.

#### GetStacktrace
```go
func GetStacktrace() bool
```
Returns true if stacktrace output is currently enabled for debug messages, either through flags, environment variables, or SetStacktrace calls.

### Example Usage

```go
package main

import (
    "github.com/Merith-TK/utils/pkg/debug"
)

func main() {
    // Enable debug mode
    debug.SetDebug(true)
    
    // Set a custom title
    debug.SetTitle("MYAPP")
    defer debug.ResetTitle()
    
    // Debug output
    debug.Print("Application starting")
    debug.Println("This includes a newline")
    
    // Enable self-destruct for testing
    debug.Suicide(30) // Exit after 30 seconds if suicide mode enabled
    
    // Check debug status
    if debug.GetDebug() {
        debug.Print("Debug mode is active")
    }
}
```

---

## Package: config

**Import**: `github.com/Merith-TK/utils/pkg/config`

### Overview

The config package provides configuration management utilities including environment variable manipulation, key replacement, and TOML configuration file handling.

### Functions

#### EnvKeyReplace
```go
func EnvKeyReplace(input string, replacements map[string]string) string
```
Replaces all `{key}` placeholders in the input string with their corresponding values from the replacements map. This is useful for templating configuration strings.

**Parameters**:
- `input`: String containing `{key}` placeholders
- `replacements`: Map of key-value pairs for replacement

**Returns**: String with all placeholders replaced

#### EnvOverride
```go
func EnvOverride(env map[string]string)
```
Sets multiple environment variables from the provided map. This is a convenience function for batch setting environment variables.

**Parameters**:
- `env`: Map of environment variable names to values

#### LoadToml
```go
func LoadToml(target interface{}, configfile string) error
```
Loads and parses a TOML configuration file into the provided struct pointer. The config file must be in valid TOML format.

**Parameters**:
- `target`: Pointer to struct that will receive the parsed configuration
- `configfile`: Path to the TOML configuration file

**Returns**: Error if file doesn't exist, can't be read, or contains invalid TOML

#### SaveToml
```go
func SaveToml(path string, cfg interface{}) error
```
Saves a struct as TOML format to the specified file path. Creates the file if it doesn't exist, overwrites if it does.

**Parameters**:
- `path`: File path where TOML will be saved
- `cfg`: Struct to serialize as TOML

**Returns**: Error if file can't be created or struct can't be serialized

### Example Usage

```go
package main

import (
    "fmt"
    "github.com/Merith-TK/utils/pkg/config"
)

type AppConfig struct {
    Name    string `toml:"name"`
    Version string `toml:"version"`
    Debug   bool   `toml:"debug"`
    Port    int    `toml:"port"`
}

func main() {
    // Key replacement
    replacements := map[string]string{
        "{USER}": "john",
        "{HOME}": "/home/john",
        "{APP}":  "myapp",
    }
    
    template := "User: {USER}, Home: {HOME}, App: {APP}"
    result := config.EnvKeyReplace(template, replacements)
    fmt.Println(result) // User: john, Home: /home/john, App: myapp
    
    // Environment override
    env := map[string]string{
        "DEBUG": "true",
        "PORT":  "8080",
        "MODE":  "production",
    }
    config.EnvOverride(env)
    
    // TOML configuration
    cfg := AppConfig{
        Name:    "MyApp",
        Version: "1.0.0",
        Debug:   true,
        Port:    8080,
    }
    
    // Save configuration
    err := config.SaveToml("app.toml", cfg)
    if err != nil {
        panic(err)
    }
    
    // Load configuration
    var loadedCfg AppConfig
    err = config.LoadToml(&loadedCfg, "app.toml")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Loaded config: %+v\n", loadedCfg)
}
```

---

## Package: archive

**Import**: `github.com/Merith-TK/utils/pkg/archive`

### Overview

The archive package provides utilities for working with archive formats, currently supporting ZIP extraction with security features and proper path handling.

### Functions

#### Unzip
```go
func Unzip(src, dest string) error
```
Extracts a ZIP archive from the source path to the destination directory. All files and folders in the archive will be extracted, preserving the directory structure. The function includes security measures to prevent path traversal attacks and ensures proper file permissions are maintained.

**Parameters**:
- `src`: Path to the ZIP archive file
- `dest`: Destination directory for extraction

**Returns**: Error if extraction fails for any reason

**Features**:
- Safe path handling to prevent directory traversal attacks
- Automatic parent directory creation
- Directory structure preservation
- File permission preservation
- Proper cleanup on errors

### Example Usage

```go
package main

import (
    "log"
    "github.com/Merith-TK/utils/pkg/archive"
)

func main() {
    // Extract ZIP archive
    err := archive.Unzip("data.zip", "/path/to/extract")
    if err != nil {
        log.Fatal("Failed to extract archive:", err)
    }
    
    log.Println("Archive extracted successfully")
}
```

---

## Package: driveutil

**Import**: `github.com/Merith-TK/utils/pkg/driveutil`

### Overview

The driveutil package provides comprehensive Windows drive management utilities including drive enumeration, metadata extraction, monitoring, and utility functions. This package is Windows-specific and uses the Windows API.

### Types

#### DriveInfo
```go
type DriveInfo struct {
    Letter string  // Drive letter (e.g., "C:\\")
    Label  string  // Volume label
    Serial uint32  // Volume serial number
    Type   uint32  // Drive type (DRIVE_FIXED, DRIVE_REMOVABLE, etc.)
}
```
Represents comprehensive information about a Windows drive.

#### DriveStore
```go
type DriveStore map[string]bool
```
Tracks drives by their unique identifier (combination of drive letter and serial number) to detect insertions and removals during monitoring.

### Functions

#### ListDrives
```go
func ListDrives() []DriveInfo
```
Returns a slice of DriveInfo for all present fixed and removable drives. Only includes drives that are currently accessible and have valid serial numbers.

**Returns**: Slice of DriveInfo structs for all available drives

#### GetVolumeSerialNumber
```go
func GetVolumeSerialNumber(root string) (uint32, error)
```
Returns the volume serial number for a given drive root path.

**Parameters**:
- `root`: Drive root path (e.g., "C:\\")

**Returns**: 
- `uint32`: Volume serial number
- `error`: Error if drive is not accessible or doesn't exist

#### DriveExists
```go
func DriveExists(drive string) bool
```
Checks if a drive path exists and is accessible.

**Parameters**:
- `drive`: Drive path to check (e.g., "E:\\")

**Returns**: True if drive exists and is accessible

#### (DriveStore) DetectDrives
```go
func (store DriveStore) DetectDrives(onNewDrive func(drive string, serial uint32))
```
Enumerates all logical drives, checks their type, and gets their serial number. Calls the provided callback for each new drive detected that wasn't previously in the store.

**Parameters**:
- `onNewDrive`: Callback function called for each newly detected drive

#### (DriveStore) MonitorDrives
```go
func (store DriveStore) MonitorDrives(onNewDrive func(drive string, serial uint32), interval time.Duration)
```
Continuously monitors for drive changes by calling DetectDrives in a loop with the specified interval. This function blocks and should typically be run in a goroutine.

**Parameters**:
- `onNewDrive`: Callback function called for each newly detected drive
- `interval`: Time interval between drive detection cycles

### Example Usage

```go
package main

import (
    "fmt"
    "time"
    "github.com/Merith-TK/utils/pkg/driveutil"
)

func main() {
    // List all current drives
    drives := driveutil.ListDrives()
    fmt.Println("Current drives:")
    for _, drive := range drives {
        fmt.Printf("  %s - %s (Serial: %08X, Type: %d)\n", 
            drive.Letter, drive.Label, drive.Serial, drive.Type)
    }
    
    // Check if specific drive exists
    if driveutil.DriveExists("E:\\") {
        fmt.Println("Drive E: is available")
        
        // Get its serial number
        if serial, err := driveutil.GetVolumeSerialNumber("E:\\"); err == nil {
            fmt.Printf("Drive E: serial number: %08X\n", serial)
        }
    }
    
    // Monitor for new drives
    store := make(driveutil.DriveStore)
    
    // Initial detection
    store.DetectDrives(func(drive string, serial uint32) {
        fmt.Printf("Initial drive found: %s (Serial: %08X)\n", drive, serial)
    })
    
    // Continuous monitoring (run in goroutine for non-blocking)
    go store.MonitorDrives(func(drive string, serial uint32) {
        fmt.Printf("New drive detected: %s (Serial: %08X)\n", drive, serial)
    }, 2*time.Second)
    
    // Keep main goroutine alive
    select {}
}
```

### Windows Drive Types

The `Type` field in `DriveInfo` corresponds to Windows drive types:

| Constant | Value | Description |
|----------|-------|-------------|
| `DRIVE_UNKNOWN` | 0 | Drive type unknown |
| `DRIVE_NO_ROOT_DIR` | 1 | Invalid root path |
| `DRIVE_REMOVABLE` | 2 | Removable drive (floppy, USB, etc.) |
| `DRIVE_FIXED` | 3 | Fixed drive (hard disk) |
| `DRIVE_REMOTE` | 4 | Network drive |
| `DRIVE_CDROM` | 5 | CD-ROM drive |
| `DRIVE_RAMDISK` | 6 | RAM disk |

The driveutil package only monitors `DRIVE_REMOVABLE` and `DRIVE_FIXED` drives by default.

---

## Error Handling

All packages follow Go's standard error handling conventions:

- Functions that can fail return an `error` as their last return value
- `nil` error indicates success
- Non-nil error contains descriptive error message
- Errors should be checked and handled appropriately

## Thread Safety

- **debug package**: Not thread-safe. Global state modifications should be synchronized if used across goroutines
- **config package**: Thread-safe for read operations, not thread-safe for concurrent file operations
- **archive package**: Thread-safe, no shared state
- **driveutil package**: Thread-safe for individual function calls, DriveStore should be protected if accessed from multiple goroutines

## Platform Compatibility

- **debug package**: Cross-platform
- **config package**: Cross-platform  
- **archive package**: Cross-platform
- **driveutil package**: Windows only (uses Windows API)

---

*This API reference covers all public functions and types in the utils package suite. For implementation details and examples, refer to the source code and comprehensive documentation.*