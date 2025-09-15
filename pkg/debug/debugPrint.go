// Package debug provides comprehensive debugging utilities including conditional logging,
// stacktrace output, and self-destruct functionality for development and testing.
//
// The package supports multiple debugging modes:
//   - Debug mode: Enables debug output with optional custom titles
//   - Stacktrace mode: Includes filtered stack traces in debug output
//   - Suicide mode: Allows processes to self-terminate after a timeout
//
// Configuration can be done via command-line flags or environment variables:
//   - -debug or DEBUG=true: Enable debug output
//   - -stacktrace or STACKTRACE=true: Enable stacktrace output
//   - -suicide or SUICIDE=true: Enable suicide mode
//
// Example usage:
//
//	debug.SetDebug(true)
//	debug.SetTitle("MYAPP")
//	debug.Print("This is a debug message")
//	debug.Suicide(30) // Self-destruct after 30 seconds if suicide mode enabled
package debug

import (
	"flag"
	"fmt"
	"log"
	"os"
	runDebug "runtime/debug"
	"strings"
	"time"
)

// defaultTitle is the default prefix for debug messages.
const defaultTitle = ""

var (
	// enableDebug indicates if debug mode is enabled.
	enableDebug bool = false
	// enableStacktrace indicates if stacktrace output is enabled.
	enableStacktrace bool = false
	// enableSuicide indicates if suicide mode is enabled.
	enableSuicide bool = false
	// Title is the current debug title prefix.
	Title string = defaultTitle
)

// init registers debug, stacktrace, and suicide flags.
func init() {
	// register debug flag
	if flag.Lookup("debug") == nil {
		flag.BoolVar(&enableDebug, "debug", flag.Lookup("debug") != nil || os.Getenv("DEBUG") == "true", "Enable Debug Mode")
	}
	if flag.Lookup("stacktrace") == nil {
		flag.BoolVar(&enableStacktrace, "stacktrace", flag.Lookup("stacktrace") != nil || os.Getenv("STACKTRACE") == "true", "Enable Stacktrace")
	}
	if flag.Lookup("suicide") == nil {
		flag.BoolVar(&enableSuicide, "suicide", flag.Lookup("suicide") != nil || os.Getenv("SUICIDE") == "true", "Enable Suicide Mode")
	}
}

// Suicide enables a self-destruct mechanism that will terminate the process after the specified
// timeout in seconds, but only if suicide mode is enabled via flag or environment variable.
// This is primarily used for testing and development to prevent runaway processes.
// The function is non-blocking and starts a goroutine to handle the timeout.
func Suicide(timeout int) {
	if enableSuicide {
		go func() {
			time.Sleep(time.Duration(timeout) * time.Second)
			log.Printf("[TIMEOUT] Exiting after %d seconds (self-destruct)", timeout)
			os.Exit(0)
		}()
	}
}

// Println prints the given message to the standard output followed by a newline character,
// but only if debug mode is enabled. This is a convenience wrapper around Print that
// automatically adds a newline. The message is prefixed with [DEBUG] and optional title.
func Println(message ...any) {
	Print(message)
	fmt.Print("\n")
}

// Print outputs the given message to standard output if debug mode is enabled.
// Messages are prefixed with [DEBUG] and optionally include a custom title if set.
// If stacktrace mode is also enabled, a filtered stack trace is included that
// excludes internal debug package frames for cleaner output.
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
			if !strings.Contains(lines[i], "github.com/Merith-TK/utils/pkg/debug") && !strings.Contains(lines[i], "runtime/debug") {
				newLines = append(newLines, lines[i])
			}
		}
		fmt.Print(strings.Join(newLines, "\n"))
	}
}
