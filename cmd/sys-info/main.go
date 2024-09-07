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
