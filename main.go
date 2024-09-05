package main

import (
	"flag"

	"github.com/Merith-TK/utils/debug"
)

func main() {
	flag.Parse()

	// this will be expanded as more and more tests are added
	if !debug.Enabled {
		flag.Usage()
	}
	debug.Title = "Testing"
	debug.Print("Hello", "World")
}
