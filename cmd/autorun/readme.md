# Autorun Drive Manager

A secure autorun manager for Windows that provides sandboxed execution of programs from removable drives with a modern GUI interface and security prompts.

## Features

- **Modern GUI**: Clean, card-based interface for managing drive configurations
- **Security Prompts**: Alerts when unknown autorun configurations are detected
- **Sandboxed Execution**: Advanced Windows isolation for secure program execution
- **Environment Isolation**: Custom environment variables and directory redirection
- **Drive Monitoring**: Real-time detection of new drives and autorun configurations
- **Configuration Management**: Easy-to-use TOML-based configuration system
- **System Tray Integration**: Runs quietly in background with tray access

## Installation

```bash
go build -o autorun.exe
# Or use the install flag to add to Windows startup
autorun.exe -install
```

## Usage

### Basic Usage
```bash
# Run the autorun manager
autorun.exe

# Install to Windows startup folder
autorun.exe -install

# Run with timeout (for testing)
autorun.exe -timeout 60
```

### Configuration File Format

Autorun configurations are stored as `.autorun.toml` files in the root of drives:

```toml
autorun = "/setup.exe"
workDir = "/installer"
isolated = true

[environment]
LANG = "en_US"
CUSTOM_VAR = "value"
```

### Security Features

When a drive with an unknown autorun configuration is detected, the security dialog shows:
- Drive information and config hash (MD5)
- Configuration details (command, working directory, isolation status)
- Environment variables
- User choice: Allow, Allow Once, Deny, Deny Once

Security decisions are stored by drive serial number, not drive letter.

### Isolation Mode

When isolation is enabled:
- **Filesystem restrictions**: Limited to own drive only, cannot access C:\ or other drives
- **Environment isolation**: Redirected user directories (AppData, Temp, etc.)
- **Process limits**: 512MB memory limit and 5-minute timeout
- **Job isolation**: All child processes contained within sandbox

## Files

- `main.go` - Application entry point and initialization
- `ui.go` - Main GUI interface with drive listing
- `dialog.go` - Configuration dialog for editing autorun settings
- `security.go` - Security metadata management and decision storage
- `security_dialog.go` - Security prompt dialog for unknown configurations
- `monitor.go` - Drive monitoring and autorun execution logic
- `autorun.go` - Core autorun execution with sandboxing
- `sandbox_windows.go` - Windows-specific sandboxing implementation
- `tray.go` - System tray functionality
- `installer.go` - Installation to Windows startup
- `types.go` - Type definitions and global variables

## Dependencies

- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit
- [systray](https://github.com/getlantern/systray) - System tray integration
- [toml](https://github.com/BurntSushi/toml) - TOML configuration parsing
- [windows](https://golang.org/x/sys/windows) - Windows API access

## Security Considerations

This application uses Windows job objects and restricted process creation for sandboxing. Some features may require administrator privileges for full security isolation. The application gracefully falls back to environment-only isolation when advanced sandboxing is not available.
    FOO = "BAR"
```

### Configuration Options:
- `program`: The program to run (required).
- `workDir`: Optional. The working directory for the program (defaults to USB root).
- `isolated`: Optional. If true, clears the system environment variables before running the program, ensuring no external variables interfere.
- `environment`: Optional. Define custom key-value pairs that will be added as environment variables for the program.

### Placeholder Support:
The configuration supports two placeholder values for dynamic paths:
- `{drive}`: Refers to the root of the USB drive.
- `{work}`: Refers to the working directory specified in `workDir`.

You can use these placeholders in your `.autorun.toml` file for flexible path handling.

For example:
```toml
program = "{work}/my_program.exe"
workDir = "{drive}/scripts"
```

This configuration will run `my_program.exe` from the `/scripts` folder on the USB drive.

## Contributing
Contributions are welcome! Please submit issues or pull requests via GitHub.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
