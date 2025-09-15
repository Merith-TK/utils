// Package main implements a system information display utility that shows
// comprehensive details about the operating system, hardware, and Go runtime.
//
// The sys-info utility provides essential system information including:
//   - Operating system and architecture detection
//   - CPU core count information
//   - Go compiler and runtime details
//   - Build target information
//   - Go build information and module details
//
// Usage:
//   sys-info
//
// Output includes:
//   - Operating System (e.g., windows, linux, darwin)
//   - Architecture (e.g., amd64, arm64, 386)
//   - Number of CPU cores available
//   - Go compiler version and type
//   - Build target platform
//   - Detailed build information from debug.ReadBuildInfo()
//
// This tool is useful for system diagnostics, environment verification,
// and debugging deployment issues across different platforms.
package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

func main() {
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Number of CPUs:", runtime.NumCPU())

	fmt.Println("Compiler:", runtime.Compiler)
	fmt.Println("Build Target:", runtime.GOOS+"/"+runtime.GOARCH)
	fmt.Println("")

	buildinfo, ok := debug.ReadBuildInfo()
	fmt.Println("Build Info:", buildinfo, ok)

}
