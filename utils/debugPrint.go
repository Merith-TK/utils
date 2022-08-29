package utils

import "log"

var DebugMode bool = false

// DebugPrint prints the given message to the log if the debug flag is set.
func DebugPrint(message ...any) {
	if DebugMode {
		log.Println(message...)
	}
}
