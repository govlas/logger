// Package logger provides collection of loggin functions with colored output and stack frames
package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/shiena/ansicolor"
)

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

	output     = ansicolor.NewAnsiColorWriter(os.Stdout)
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

// SetOut set writer fo print logs. Default writer is os.Stdout
func SetOut(out io.Writer) {
	output = ansicolor.NewAnsiColorWriter(out)
}

func makeStack() (ret string) {
	st := make([]uintptr, 50)
	cnt := runtime.Callers(4, st)
	st = st[0:cnt]
	for _, r := range st {
		fnc := runtime.FuncForPC(r)
		file, line := fnc.FileLine(r - 1)

		stack_line := fmt.Sprintf("\t%s:%d\n\t\t%s\n", file, line, fnc.Name())
		if colored {
			stack_line = color.MagentaString(stack_line)
		}

		ret += stack_line
	}
	return
}

func printLog(colora color.Attribute, prefix string, printTrace bool, message string, v ...interface{}) {
	var f func(format string, a ...interface{}) string
	if colored {
		f = color.New(colora).SprintfFunc()
	} else {
		f = fmt.Sprintf
	}

	var fileName string
	if fname != FileNameNo {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		if fname == FileNameShort {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		fileName = fmt.Sprintf(" %s:%d", file, line)
	}

	var stack string
	if printTrace && btrace {
		stack = makeStack()
	}
	fmt.Fprint(output, f("%s %s %s\n", time.Now().Format(timeFormat), prefix+fileName, fmt.Sprintf(message, v...)), stack)
}

// Info prints info log (green if colored)
func Info(message string, v ...interface{}) {
	printLog(color.FgGreen, InfoPrefix, false, message, v...)
}

// Warning prints warning log (yellow if colored)
func Warning(message string, v ...interface{}) {
	printLog(color.FgYellow, WarnPrefix, true, message, v...)
}

// Error prints error log (red if colored)
func Error(message string, v ...interface{}) {
	printLog(color.FgRed, ErrPrefix, true, message, v...)
}

// Debug prints debug log (blue if colored)
func Debug(message string, v ...interface{}) {
	if pdebug {
		printLog(color.FgBlue, DebugPrefix, false, message, v...)
	}
}

// Fatal prints error log and finish proccess by os.Exit(1)
func Fatal(message string, v ...interface{}) {
	Error(message, v...)
	os.Exit(1)
}

// Panic prints error log and call panic
func Panic(message string, v ...interface{}) {
	DisableBTrace()
	Error(message, v...)
	panic(fmt.Sprintf(message, v...))
}

// ErrorErr prints err as error log and returns true if err!=nil
func ErrorErr(err error) bool {
	if err != nil {
		Error("%v", err)
		return true
	}
	return false
}

// WarningErr prints err as warning log and returns true if err!=nil
func WarningErr(err error) bool {
	if err != nil {
		Warning("%v", err)
		return true
	}
	return false
}

// FatalErr prints err as fatal log
func FatalErr(err error) {
	if err != nil {
		Fatal("%v", err)
	}
}
