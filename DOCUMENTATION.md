# Utils - Comprehensive Documentation

A collection of utility tools and packages for various system operations, primarily focused on Windows environments.

## Table of Contents

- [Overview](#overview)
- [Commands](#commands)
- [Packages](#packages)
- [Installation](#installation)
- [Development](#development)
- [API Reference](#api-reference)

## Overview

This repository contains a comprehensive suite of utilities organized into two main categories:

1. **Commands (`cmd/`)**: Standalone executable tools for specific tasks
2. **Packages (`pkg/`)**: Reusable Go libraries for common functionality

The project is designed with modularity in mind, allowing individual tools to be built and used independently while sharing common functionality through the package libraries.

## Commands

### autorun
**Path**: `cmd/autorun/`

A Windows system tray application that monitors removable drives and automatically executes configured programs when drives are inserted.

**Features**:
- System tray integration with Fyne UI
- Drive monitoring and detection
- Security dialog for drive access
- Configuration management per drive
- Installer functionality for startup folder

**Usage**:
```bash
autorun [-install] [-timeout seconds]
```

**Flags**:
- `-install, -i`: Install autorun service to startup folder
- `-timeout`: Exit after N seconds (for testing)

### dcomp
**Path**: `cmd/dcomp/`

A Docker Compose wrapper that provides simplified commands for common Docker Compose operations.

**Features**:
- Simplified command aliases for Docker Compose
- Direct passthrough to docker compose for unknown commands
- Streamlined workflow for development environments

**Usage**:
```bash
dcomp [command] [args...]
```

**Available Commands**:
- `build`: docker compose build
- `pull`: docker compose pull  
- `start`: docker compose up -d --remove-orphans
- `stop`: docker compose down
- `ps`: docker compose ps
- `logs`: docker compose logs

### doh-poke
**Path**: `cmd/doh-poke/`

A DNS-over-HTTPS (DoH) client that queries DNS servers using the DoH protocol.

**Features**:
- DoH query support with JSON response parsing
- Automatic HTTPS scheme detection
- Configurable DoH server endpoints
- A record resolution

**Usage**:
```bash
doh-poke [doh-server] [domain]
```

**Example**:
```bash
doh-poke cloudflare-dns.com example.com
doh-poke https://1.1.1.1/dns-query google.com
```

### doh2dns
**Path**: `cmd/doh2dns/`

A DNS-over-HTTPS to traditional DNS proxy server.

**Features**:
- Converts traditional DNS queries to DoH requests
- Configurable upstream DoH servers
- Standard DNS server interface

### downtime-client
**Path**: `cmd/downtime-client/`

Client component for a downtime monitoring system.

**Features**:
- Connects to downtime monitoring server
- Reports system status and availability
- Configurable check intervals

### downtime-server  
**Path**: `cmd/downtime-server/`

Server component for a downtime monitoring system.

**Features**:
- Receives status reports from clients
- Web interface for monitoring
- Alert system for downtime detection

### git-sort-repo
**Path**: `cmd/git-sort-repo/`

A utility for organizing and sorting Git repositories.

**Features**:
- Repository organization tools
- Batch operations on multiple repositories
- Git workflow automation

### mc-logChat
**Path**: `cmd/mc-logChat/`

Minecraft server log chat extractor and processor.

**Features**:
- Extracts chat messages from Minecraft server logs
- Filters and formats chat output
- Real-time log monitoring

### mc-server-icon
**Path**: `cmd/mc-server-icon/`

Minecraft server icon generator and manager.

**Features**:
- Server icon creation and conversion
- Batch processing of server icons
- Format validation and optimization

### mc2se
**Path**: `cmd/mc2se/`

Minecraft to Space Engineers schematic converter.

**Features**:
- Converts Minecraft structures to Space Engineers blueprints
- Block mapping and translation
- Litematica schematic support

### moba-cracker
**Path**: `cmd/moba-cracker/`

MobaXterm session file processor and password recovery tool.

**Features**:
- Session file analysis
- Password recovery utilities
- Configuration extraction

### sys-info
**Path**: `cmd/sys-info/`

System information display utility.

**Features**:
- Operating system detection
- Hardware architecture information
- Go runtime information
- Build information display

**Usage**:
```bash
sys-info
```

**Output Example**:
```
Operating System: windows
Architecture: amd64
Number of CPUs: 8
Compiler: gc
Build Target: windows/amd64
```

### testdriveutil
**Path**: `cmd/testdriveutil/`

Test utility for the driveutil package functionality.

**Features**:
- Drive detection testing
- Serial number verification
- Drive monitoring demonstration

### traytest
**Path**: `cmd/traytest/`

System tray functionality testing utility.

**Features**:
- System tray integration testing
- UI component verification
- Cross-platform tray behavior testing

### wscli
**Path**: `cmd/wscli/`

WebSocket command-line interface client.

**Features**:
- WebSocket connection management
- Interactive command-line interface
- Message sending and receiving
- Connection status monitoring

## Packages

### debug (`pkg/debug/`)

Comprehensive debugging utilities with conditional logging, stacktrace output, and self-destruct functionality.

**Key Features**:
- Conditional debug output based on flags or environment variables
- Custom title prefixes for debug messages
- Filtered stacktrace output excluding internal frames
- Self-destruct mechanism for testing (suicide mode)

**Configuration**:
- `-debug` flag or `DEBUG=true` environment variable
- `-stacktrace` flag or `STACKTRACE=true` environment variable  
- `-suicide` flag or `SUICIDE=true` environment variable

**API**:
```go
// Enable/disable debug mode
debug.SetDebug(true)
debug.GetDebug() bool

// Set custom title prefix
debug.SetTitle("MYAPP")
defer debug.ResetTitle()

// Debug output
debug.Print("Debug message")
debug.Println("Debug message with newline")

// Self-destruct after timeout
debug.Suicide(30) // 30 seconds
```

### config (`pkg/config/`)

Configuration management utilities for environment variables, key replacement, and TOML file handling.

**Key Features**:
- Template-based key replacement in strings
- Environment variable batch setting
- TOML configuration file loading and saving
- Cross-platform path handling

**API**:
```go
// Key replacement
replacements := map[string]string{"{USER}": "john", "{HOME}": "/home/john"}
result := config.EnvKeyReplace("User: {USER}, Home: {HOME}", replacements)

// Environment override
env := map[string]string{"DEBUG": "true", "PORT": "8080"}
config.EnvOverride(env)

// TOML configuration
type Config struct {
    Name    string `toml:"name"`
    Version string `toml:"version"`
}

var cfg Config
err := config.LoadToml(&cfg, "config.toml")
err = config.SaveToml("config.toml", cfg)
```

### archive (`pkg/archive/`)

Archive format utilities, currently supporting ZIP extraction with security features.

**Key Features**:
- Safe ZIP extraction with path traversal protection
- Directory structure preservation
- Automatic parent directory creation
- File permission preservation

**API**:
```go
// Extract ZIP archive
err := archive.Unzip("archive.zip", "/path/to/extract")
if err != nil {
    log.Fatal("Failed to extract:", err)
}
```

### driveutil (`pkg/driveutil/`)

Windows-specific drive management utilities for enumeration, monitoring, and metadata extraction.

**Key Features**:
- Drive detection and enumeration (fixed and removable)
- Volume serial number extraction
- Drive monitoring with callback support
- Drive existence checking
- Comprehensive drive metadata

**Types**:
```go
type DriveInfo struct {
    Letter string  // Drive letter (e.g., "C:\\")
    Label  string  // Volume label
    Serial uint32  // Volume serial number
    Type   uint32  // Drive type (DRIVE_FIXED, DRIVE_REMOVABLE, etc.)
}

type DriveStore map[string]bool // Tracks drives by unique ID
```

**API**:
```go
// List all drives
drives := driveutil.ListDrives()
for _, drive := range drives {
    fmt.Printf("Drive: %s, Label: %s, Serial: %08X\n", 
        drive.Letter, drive.Label, drive.Serial)
}

// Monitor for new drives
store := make(driveutil.DriveStore)
store.DetectDrives(func(drive string, serial uint32) {
    fmt.Printf("New drive: %s (Serial: %08X)\n", drive, serial)
})

// Continuous monitoring
go store.MonitorDrives(callback, 1*time.Second)

// Check drive existence
if driveutil.DriveExists("E:\\") {
    fmt.Println("Drive E: exists")
}

// Get volume serial number
serial, err := driveutil.GetVolumeSerialNumber("C:\\")
```

## Installation

### Prerequisites
- Go 1.23 or later
- Windows OS (for drive utilities and some commands)
- Git

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/Merith-TK/utils.git
cd utils
```

2. Build all commands:
```bash
go build -o bin/ ./cmd/...
```

3. Build specific command:
```bash
go build -o bin/sys-info.exe ./cmd/sys-info
```

4. Install as Go modules:
```bash
go install ./cmd/sys-info
```

### Using as Library

Add to your `go.mod`:
```go
require github.com/Merith-TK/utils v0.0.0
```

Import packages:
```go
import (
    "github.com/Merith-TK/utils/pkg/debug"
    "github.com/Merith-TK/utils/pkg/config"
    "github.com/Merith-TK/utils/pkg/archive"
    "github.com/Merith-TK/utils/pkg/driveutil"
)
```

## Development

### Project Structure
```
utils/
├── cmd/                    # Executable commands
│   ├── autorun/           # Drive autorun manager
│   ├── dcomp/             # Docker Compose wrapper
│   ├── doh-poke/          # DoH DNS client
│   └── ...                # Other commands
├── pkg/                   # Reusable packages
│   ├── debug/             # Debug utilities
│   ├── config/            # Configuration management
│   ├── archive/           # Archive utilities
│   └── driveutil/         # Drive utilities
├── main.go                # Test runner for all packages
├── go.mod                 # Go module definition
└── README.md              # Project overview
```

### Testing

Run the comprehensive test suite:
```bash
go run main.go
```

This will test all packages with sample data and demonstrate functionality.

### Building

Use the provided Magefile for build automation:
```bash
mage build
```

Or build manually:
```bash
# Build all commands
for dir in cmd/*/; do
    go build -o "bin/$(basename "$dir")" "./$dir"
done
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Update documentation
6. Submit a pull request

## API Reference

### Debug Package

#### Functions

- `SetDebug(enabled bool)` - Enable/disable debug mode programmatically
- `GetDebug() bool` - Check if debug mode is enabled
- `SetTitle(title string)` - Set custom debug message prefix
- `GetTitle() string` - Get current debug title
- `ResetTitle()` - Reset debug title to default
- `Print(message ...any)` - Output debug message if debug enabled
- `Println(message ...any)` - Output debug message with newline
- `Suicide(timeout int)` - Enable self-destruct after timeout seconds
- `SetStacktrace(enabled bool)` - Enable/disable stacktrace output
- `GetStacktrace() bool` - Check if stacktrace is enabled

### Config Package

#### Functions

- `EnvKeyReplace(input string, replacements map[string]string) string` - Replace {key} placeholders
- `EnvOverride(env map[string]string)` - Set multiple environment variables
- `LoadToml(target interface{}, configfile string) error` - Load TOML configuration
- `SaveToml(path string, cfg interface{}) error` - Save configuration as TOML

### Archive Package

#### Functions

- `Unzip(src, dest string) error` - Extract ZIP archive to destination

### DriveUtil Package

#### Types

```go
type DriveInfo struct {
    Letter string
    Label  string
    Serial uint32
    Type   uint32
}

type DriveStore map[string]bool
```

#### Functions

- `ListDrives() []DriveInfo` - Get all available drives
- `GetVolumeSerialNumber(root string) (uint32, error)` - Get drive serial number
- `DriveExists(drive string) bool` - Check if drive exists
- `(store DriveStore) DetectDrives(callback func(string, uint32))` - Detect new drives
- `(store DriveStore) MonitorDrives(callback func(string, uint32), interval time.Duration)` - Monitor drives continuously

---

*This documentation covers the complete utils package suite. For specific implementation details, refer to the source code and inline documentation.*