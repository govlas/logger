// Package logger provides collection of loggin functions with colored output and stack frames
package logger

const (
	// FileNameNo indicates to no print filename
	FileNameNo int = iota
	// FileNameShort indicates to print short filename
	FileNameShort
	// FileNameLong indicates to print full filename
	FileNameLong
)

var (
	colored = false
	btrace  = true
	pdebug  = true
	fname   = FileNameNo

	timeFormat = "2006-Jan-2 15:04:05.0000"

	// InfoPrefix is text prefix for info log
	InfoPrefix = "INFO"
	// WarnPrefix is text prefix for warning log
	WarnPrefix = "WARN"
	// ErrPrefix is text prefix for error log
	ErrPrefix = " ERR"
	// DebugPrefix is text prefix for debug log
	DebugPrefix = " DEB"
)

// EnableColored enables colored output
func EnableColored() {
	colored = true
}

// DisableColored disables colored output
func DisableColored() {
	colored = false
}

// EnableBTrace enables backtrace for warnings and errors
func EnableBTrace() {
	btrace = true
}

// DisableBTrace disables backtrace
func DisableBTrace() {
	btrace = false
}

// EnableDebug enables debug logs
func EnableDebug() {
	pdebug = true
}

// DisableDebug disables debug logs
func DisableDebug() {
	pdebug = false
}

// SetFileName sets flag for print ceurrent file name in log
func SetFileName(fn int) {
	switch fn {
	case FileNameLong, FileNameNo, FileNameShort:
		fname = fn
	}
}
