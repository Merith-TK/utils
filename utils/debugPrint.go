package utils

import "log"

var (
	DebugMode  bool   = false
	DebugTitle string = "DEBUG"
)

// DebugPrint prints the given message to the log if the debug flag is set.
func DebugPrint(message ...any) {
	if DebugMode {
		// add DebugTitle to the message
		log.Println(DebugTitle, message)
	}
}
