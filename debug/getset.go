package debug

func SetTitle(title string) {
	Title = title
}
func GetTitle() string {
	return Title
}
func ResetTitle() {
	Title = defaultTitle
}

func SetDebug(enabled bool) {
	enableDebug = enabled
}
func GetDebug() bool {
	return enableDebug
}

func SetStacktrace(enabled bool) {
	enableStacktrace = enabled
}
func GetStacktrace() bool {
	return enableStacktrace
}
