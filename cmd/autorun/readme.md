# autorun

This program is a reimplementation of the "autorun" feature that was removed from Windows.

## Description
The autorun program allows you to automatically run a specified executable or script when a USB drive is inserted into your computer. It provides a convenient way to automate tasks or launch applications upon USB drive detection.

## Features
- Automatically execute a specified file when a USB drive is connected
- Customizable file execution options
- Cross-platform compatibility (Windows, macOS, Linux)

## Installation
To install `autorun`, you can use the following command:

```shell
go get github.com/merith-tk/utils/cmd/autorun@latest
```

This will download and install the `autorun` package from the GitHub repository `github.com/merith-tk/autorun`.

Once the installation is complete, you can use `autorun install` to enable autostart on login

## Usage
1. The USB drives should have a `.autorun.toml` file in the root directory.
2. Run the `autorun` executable.
3. Insert a USB drive into your computer and watch the magic happen!


## Configuration

To configure the `autorun` program, you can create a `.autorun.toml` file in the root directory of your USB drives. Here is an example configuration section for the file:

```toml
autorun = "example.exe"
workDir = "./"
isolated = false

[environment]
    FOO = "BAR"
```

In this configuration section:
- `program` specifies the program to run, along with any arguments.
- `workDir` specifies the relative path to be used as the base directory. If not found, the default is the root directory of the USB drive.
- `isolated` determines whether the environment should be emptied before running the program, allowing for true portability.
- `environment` is a section where you can define key-value pairs to inject into the program's environment.

Feel free to modify the values according to your needs.
To configure the `autorun` program, you can create a `.autorun.toml` file in the root directory of your USB drives. Here is an example configuration section for the file:

```toml
program = "example.exe"
# Optional: workDir specifies the relative path to be used as the base directory. If not found, the default is the root directory of the USB drive.
workDir = "./"
# Optional: isolated determines whether the environment should be emptied before running the program, allowing for true portability.
isolated = false

[environment]
    # Optional: Define key-value pairs to inject into the program's environment.
    # FOO = "BAR"
```

In this configuration section:
- `program` specifies the program to run, along with any arguments.
- `workDir` specifies the relative path to be used as the base directory. If not found, the default is the root directory of the USB drive.
- `isolated` determines whether the environment should be emptied before running the program, allowing for true portability.
- `environment` is a section where you can define key-value pairs to inject into the program's environment.

Please note that only `program` is required, while the rest are optional. Feel free to modify the values according to your needs.


## Contributing
Contributions are welcome! If you have any ideas, suggestions, or bug reports, please open an issue or submit a pull request.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
