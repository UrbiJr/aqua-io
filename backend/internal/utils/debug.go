package utils

// DebugEnabled when true, will enable debug logs and other debug related functions.
var DebugEnabled bool

// SetDebug is a convenience function used to switch the DebugEnabled boolean.
func SetDebug(v bool) {
	DebugEnabled = v
}
