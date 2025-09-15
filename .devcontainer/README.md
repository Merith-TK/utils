# Dev Container Configuration

Simple dev container setup for Go development using the official Microsoft Go dev container image.

## Features

- **Go 1.23**: Latest stable Go version with all standard tools
- **Docker-in-Docker**: For containerized development workflows
- **golang.go**: Official Go VS Code extension with default settings

## Usage

The dev container uses the official `mcr.microsoft.com/devcontainers/go:1.23` image which includes:
- Go 1.23 with all standard tools (go, gofmt, etc.)
- Common development tools (git, curl, wget, etc.)
- Node.js and npm for web development
- Default Go development environment

### Building the Project
```bash
# Download dependencies
go mod download

# Build all commands
go build -o bin/ ./cmd/...

# Build specific command
go build -o bin/sys-info ./cmd/sys-info

# Run tests
go test ./...

# Run with mage (if available)
mage build
```

### Go Tools
The Go extension will automatically prompt to install additional tools like:
- gopls (language server)
- dlv (debugger)
- staticcheck (linter)
- And others as needed

## Rebuilding the Container

To rebuild the dev container:
1. Open Command Palette (Ctrl+Shift+P)
2. Run "Dev Containers: Rebuild Container"

## Configuration

The setup uses minimal configuration to rely on sensible defaults:
- Official Microsoft Go dev container image
- Docker-in-Docker feature for containerized workflows
- Only the essential golang.go extension

All other settings use VS Code and Go extension defaults.