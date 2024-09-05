package debug

import (
	"flag"
	"log"
	"os"

	runDebug "runtime/debug"
)

var (
	Enabled bool   = false
	Title   string = "DEBUG"
)

func init() {
	// register debug flag
	flag.BoolVar(&Enabled, "debug", flag.Lookup("debug") != nil || os.Getenv("DEBUG") == "true", "Enable Debug Mode")
}

// Usage:
//		debug.Enabled = true
// 	    debug.Title = "HELLO"
//		debug.Print("Hello World")
// Output: DEBUG Hello World

func Print(message ...any) {
	log.Println("[DEBUG]", Title, message, "\n"+string(runDebug.Stack()))
}
