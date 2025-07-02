package debug

// Sets the title of the debug message.
// note this currently is set globally and will affect all debug messages.
// it is reccomended to use this function at the start of a function followed by `defer debug.ResetTitle()`
func SetTitle(title string) {
	Title = title
}

// Gets the current title of the debug message.
func GetTitle() string {
	return Title
}

// Resets the title of the debug message to the default value.
func ResetTitle() {
	Title = defaultTitle
}

// Toggle the debug regardless of the flag.
func SetDebug(enabled bool) {
	enableDebug = enabled
}

// Get the current debug status.
func GetDebug() bool {
	return enableDebug
}

// Toggle the stacktrace regardless of the flag.
func SetStacktrace(enabled bool) {
	enableStacktrace = enabled
}

// Get the current stacktrace status.
func GetStacktrace() bool {
	return enableStacktrace
}
