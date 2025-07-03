# debug

This package provides utilities for debugging with conditional output and stacktraces.

## Variables

- `Title string` - Current debug title prefix

## Functions

- `Print(message ...any)` - Prints debug message if debug mode enabled
- `Println(message ...any)` - Prints debug message with newline
- `SetTitle(title string)` - Sets debug message title prefix
- `GetTitle() string` - Gets current debug title
- `ResetTitle()` - Resets title to default
- `SetDebug(enabled bool)` - Toggles debug mode
- `GetDebug() bool` - Gets debug mode status
- `SetStacktrace(enabled bool)` - Toggles stacktrace mode
- `GetStacktrace() bool` - Gets stacktrace mode status
- `Suicide(timeout int)` - Self-destructs after timeout if suicide mode enabled

## Flags

- `-debug` - Enable debug mode (or set DEBUG=true env var)
- `-stacktrace` - Enable stacktraces (or set STACKTRACE=true env var)
- `-suicide` - Enable suicide mode (or set SUICIDE=true env var)

## Example

```go
import "github.com/Merith-TK/utils/pkg/debug"

debug.SetDebug(true)
debug.SetTitle("MyApp")
debug.Print("Starting application")
debug.Println("Configuration loaded")
defer debug.ResetTitle()
``` 