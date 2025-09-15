package debug

// SetTitle sets the title of the debug message.
// Note: this is set globally and will affect all debug messages.
// It is recommended to use this function at the start of a function followed by defer debug.ResetTitle().
func SetTitle(title string) {
	Title = title
}

// GetTitle returns the currently set debug message title prefix.
// Returns an empty string if no custom title has been set.
func GetTitle() string {
	return Title
}

// ResetTitle resets the debug message title prefix to the default empty value.
// This should be called to clean up after using SetTitle, typically with defer.
func ResetTitle() {
	Title = defaultTitle
}

// SetDebug programmatically enables or disables debug mode, overriding
// any command-line flag or environment variable settings.
func SetDebug(enabled bool) {
	enableDebug = enabled
}

// GetDebug returns true if debug mode is currently enabled,
// either through flags, environment variables, or SetDebug calls.
func GetDebug() bool {
	return enableDebug
}

// SetStacktrace programmatically enables or disables stacktrace output in debug messages,
// overriding any command-line flag or environment variable settings.
func SetStacktrace(enabled bool) {
	enableStacktrace = enabled
}

// GetStacktrace returns true if stacktrace output is currently enabled
// for debug messages, either through flags, environment variables, or SetStacktrace calls.
func GetStacktrace() bool {
	return enableStacktrace
}
