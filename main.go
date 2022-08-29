package main

import (
	"flag"

	"git.merith.xyz/packages/utils/debug"
)

func main() {
	flag.Parse()

	// this will be expanded as more and more tests are added
	if !debug.Enabled {
		flag.Usage()
	}

	// Test debug.Print()
	debug.Print("Hello World")
}
