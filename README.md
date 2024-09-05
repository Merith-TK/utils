# Merith's Personal Collection of Utilities and Commands

This is my personal collection of Go packages.

`main.go` is my jank way of making tests for my libraries since I have zero clue how to use Go tests.

## Installation

### Using the Script

If you want to install all utilities in this collection at once, you can use my installation script:

```bash
curl -sSL https://raw.githubusercontent.com/Merith-TK/utils/main/install.sh | bash
```

This will download and run the script, installing all the packages included in this repository.

### Installing Individual Go Packages

You can also install specific packages from this repository manually. Each package is located in the `cmd` directory. Here's how you can install them individually using Go:

1. Make sure you have Go installed. If not, download it from [here](https://golang.org/dl/).
2. Run the following command to install a specific package:

```bash
go install github.com/Merith-TK/utils/cmd/<package-name>@latest
```

Replace `<package-name>` with the actual name of the package you want to install. The package binary will be placed in your Go `bin` directory, which should be in your system's `PATH` for global use.

Ensure that your `$GOPATH/bin` is in your `PATH` so you can run the installed binaries from anywhere in your terminal.

## Notes

- `main.go` is mainly for my own testing purposes, so it's not necessarily meant for production use.
  
- Feel free to clone or fork the repository and modify it according to your needs.

---

If you have any issues or suggestions, feel free to open an issue or contribute via pull requests.
