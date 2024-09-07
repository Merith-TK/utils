package main

import (
	"flag"

	"github.com/Merith-TK/utils/debug"
)

func main() {
	flag.Parse()
	if !debug.Enabled && !debug.EnableStacktrace {
		flag.Usage()
	}
	debug.SetTitle("Testing")
	debug.Print("Hello", "World")
	debug.ResetTitle()
	debug.Print("Hello", "World")

	debug.PrintStacktrace("Hello", "Stacktrace")
}
