package main

import (
	"flag"

	"github.com/Merith-TK/utils/debug"
)

func main() {
	flag.Parse()
	if !debug.GetDebug() && !debug.GetStacktrace() {
		flag.Usage()
	}
	debug.SetTitle("Testing")
	debug.Println("Hello", "World")
	debug.ResetTitle()
	debug.Println("Hello", "World")

	main2()
}
