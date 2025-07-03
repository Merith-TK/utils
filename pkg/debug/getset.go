package debug

// SetTitle sets the title of the debug message.
// Note: this is set globally and will affect all debug messages.
// It is recommended to use this function at the start of a function followed by defer debug.ResetTitle().
func SetTitle(title string) {
	Title = title
}

// GetTitle returns the current title of the debug message.
func GetTitle() string {
	return Title
}

// ResetTitle resets the title of the debug message to the default value.
func ResetTitle() {
	Title = defaultTitle
}

// SetDebug toggles debug mode regardless of the flag.
func SetDebug(enabled bool) {
	enableDebug = enabled
}

// GetDebug returns the current debug status.
func GetDebug() bool {
	return enableDebug
}

// SetStacktrace toggles stacktrace mode regardless of the flag.
func SetStacktrace(enabled bool) {
	enableStacktrace = enabled
}

// GetStacktrace returns the current stacktrace status.
func GetStacktrace() bool {
	return enableStacktrace
}
