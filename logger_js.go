// +build js

package logger

import (
	"encoding/json"
	"fmt"
	"time"

	"honnef.co/go/js/console"
)

func printLog(prefix string, printTrace bool, message string, v ...interface{}) {

	mess := fmt.Sprintf("%s %s %s", time.Now().Format(timeFormat), prefix, fmt.Sprintf(message, v...))

	switch prefix {
	case WarnPrefix:
		console.Warn(mess)
	case ErrPrefix:
		console.Error(mess)
	default:
		console.Log(mess)
	}
}

// Info prints info log (green if colored)
func Info(message string, v ...interface{}) {
	printLog(InfoPrefix, false, message, v...)
}

// Warning prints warning log (yellow if colored)
func Warning(message string, v ...interface{}) {
	printLog(WarnPrefix, true, message, v...)
}

// Error prints error log (red if colored)
func Error(message string, v ...interface{}) {
	printLog(ErrPrefix, true, message, v...)
}

// Debug prints debug log (blue if colored)
func Debug(message string, v ...interface{}) {
	if pdebug {
		printLog(DebugPrefix, false, message, v...)
	}
}

// Fatal prints error log and finish proccess by os.Exit(1)
func Fatal(message string, v ...interface{}) {
	printLog(ErrPrefix, true, message, v...)
}

// Panic prints error log and call panic
func Panic(message string, v ...interface{}) {
	DisableBTrace()
	printLog(ErrPrefix, true, message, v...)
	panic(fmt.Sprintf(message, v...))
}

// PanicRecover recovers panic an print error message
func PanicRecover() {
	if err := recover(); err != nil {
		Error("PANIC: %s", err)
	}
}

// ErrorErr prints err as error log and returns true if err!=nil
func ErrorErr(err error) bool {
	if err != nil {
		printLog(ErrPrefix, true, "%v", err)
		return true
	}
	return false
}

// WarningErr prints err as warning log and returns true if err!=nil
func WarningErr(err error) bool {
	if err != nil {
		printLog(WarnPrefix, true, "%v", err)
		return true
	}
	return false
}

// FatalErr prints err as fatal log
func FatalErr(err error) {
	if err != nil {
		printLog(ErrPrefix, true, "%v", err)
	}
}

// JSONDebug prints object in json format
func JSONDebug(a interface{}) {
	buf, err := json.MarshalIndent(a, "", "  ")
	if !WarningErr(err) {
		Debug("%s", buf)
	}
}
