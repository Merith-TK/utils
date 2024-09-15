// Package debug provides utilities for debugging purposes.
package debug

import (
	"flag"
	"fmt"
	"log"
	"os"
	runDebug "runtime/debug"
	"strings"
)

const defaultTitle = ""

var (
	enableDebug      bool   = false
	enableStacktrace bool   = false
	Title            string = defaultTitle
)

func init() {
	// register debug flag
	if flag.Lookup("debug") == nil {
		flag.BoolVar(&enableDebug, "debug", flag.Lookup("debug") != nil || os.Getenv("DEBUG") == "true", "Enable Debug Mode")
	}
	if flag.Lookup("stacktrace") == nil {
		flag.BoolVar(&enableStacktrace, "stacktrace", flag.Lookup("stacktrace") != nil || os.Getenv("STACKTRACE") == "true", "Enable Stacktrace")
	}
}

// Println prints the given message to the standard output, followed by a newline character.
// Note: Println should not be necessary, but it is included in case it is needed.
func Println(message ...any) {
	Print(message)
	fmt.Print("\n")
}

// Print prints the given message to the standard output.
func Print(message ...any) {
	if enableDebug {
		if Title == defaultTitle {
			log.Print("[DEBUG]", message)
		} else {
			log.Print("[DEBUG] {"+Title+"}", message)
		}
	}
	if enableStacktrace && enableDebug {
		stack := runDebug.Stack()
		lines := strings.Split(string(stack), "\n")
		var newLines []string
		for i := 0; i < len(lines); i++ {
			if !strings.Contains(lines[i], "github.com/Merith-TK/utils/debug") && !strings.Contains(lines[i], "runtime/debug") {
				newLines = append(newLines, lines[i])
			}
		}
		fmt.Print(strings.Join(newLines, "\n"))
	}
}
