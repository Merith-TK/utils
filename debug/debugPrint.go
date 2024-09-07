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
	Enabled          bool   = false
	EnableStacktrace bool   = false
	Title            string = defaultTitle
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

func SetTitle(title string) {
	Title = title
}
func ResetTitle() {
	Title = defaultTitle
}

func Print(message ...any) {
	if Enabled {
		if Title == defaultTitle {
			log.Println("[DEBUG]", message)
		} else {
			log.Println("[DEBUG] {"+Title+"}", message)
		}
	}
	if Enabled && EnableStacktrace {
		stack := runDebug.Stack()
		lines := strings.Split(string(stack), "\n")
		if len(lines) > 2 {
			line := lines[6]
			line = strings.TrimSpace(line)
			path := strings.Split(line, " ")[0]
			if !strings.Contains(path, "github.com/Merith-TK/utils/debug/debugPrint.go") {
				fmt.Println("\t", path)
			}
		}
	}
}

func PrintStacktrace(message ...any) {
	Print(message)
	stack := runDebug.Stack()
	lines := strings.Split(string(stack), "\n")
	fmt.Println(strings.Join(lines[5:], "\n"))
}
