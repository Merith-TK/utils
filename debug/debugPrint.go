package debug

import (
	"flag"
	"fmt"
	"log"
	"os"
	runDebug "runtime/debug"
)

var (
	Enabled          bool   = false
	EnableStacktrace bool   = false
	Title            string = "DEBUG"
)

func init() {
	// register debug flag
	flag.BoolVar(&Enabled, "debug", flag.Lookup("debug") != nil || os.Getenv("DEBUG") == "true", "Enable Debug Mode")
	flag.BoolVar(&EnableStacktrace, "stacktrace", flag.Lookup("stacktrace") != nil || os.Getenv("STACKTRACE") == "true", "Enable Stacktrace")
}

// Usage:
//		debug.Enabled = true
// 	    debug.Title = "HELLO"
//		debug.Print("Hello World")
// Output: DEBUG Hello World

func Print(message ...any) {
	if Enabled {
		log.Println("[DEBUG]", Title, message)
	}
	if EnableStacktrace {
		log.Println("[DEBUG]", Title+" Stacktrace")
		fmt.Println(string(runDebug.Stack()))
	}
}
