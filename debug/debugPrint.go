package debug

import (
	"flag"
	"log"
)

var (
	Enabled    bool   = false
	DebugTitle string = "DEBUG"
)

func init() {
	// register debug flag
	flag.BoolVar(&Enabled, "debug", false, "Enable Debug Mode")
}

// Usage:
//		utils.Debug.Enabled = true
//		utils.Debug.Print("Hello World")
// Output: DEBUG Hello World

func Print(message ...any) {
	if Enabled {
		log.Println(DebugTitle, message)
	}
}
