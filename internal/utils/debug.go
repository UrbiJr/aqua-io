package utils

import (
	"fmt"
	"os"
)

// DebugEnabled when true, will enable debug logs and other debug related functions.
var DebugEnabled bool

// SetDebug is a convenience function used to switch the DebugEnabled boolean.
func SetDebug(v bool) {
	DebugEnabled = v
}

// Debug is a Printf function that writes to stderr.
// This function is a noop when DebugEnabled is false.
// Usage example: utils.Debug("error occurred: %v", err)
func Debug(format string, args ...interface{}) {
	if DebugEnabled {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
