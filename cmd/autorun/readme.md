# Autorun

A reimplementation of the "autorun" feature that was removed from Windows.

## Overview
`Autorun` enables automatic execution of a specified program or script when a USB drive is inserted into a Windows computer. This offers a simple way to automate tasks or launch applications upon USB detection.

## Features
- Automatically executes a specified file when a USB drive is connected.
- Customizable execution settings.
- **Windows-only support**.

## Installation
Install `autorun` with:

```shell
go install -ldflags -H=windowsgui github.com/merith-tk/utils/cmd/autorun@latest
```

After installation, enable autorun on login with:

```shell
autorun install
```

## Usage
1. Ensure your USB drive has an `.autorun.toml` file in the root directory. If one doesn't exist, `autorun` will automatically place a generic `.autorun.toml` file on the drive that does nothing by default:
    ```toml
    program = "example.exe"
    workDir = "./"
    
    [environment]
      FOO = "BAR"
    ```
2. Run the `autorun` executable.
3. Insert your USB drive and let `autorun` handle the rest.

## Standalone Mode
In standalone mode, `autorun` uses a `.autorun.toml` file placed next to the executable instead of scanning all connected drives. This is ideal for single-use configurations.

## Configuration

To customize the behavior of `autorun`, create an `.autorun.toml` file on your USB drive. Example configuration:

```toml
program = "example.exe"
workDir = "./"
isolated = false

[environment]
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
