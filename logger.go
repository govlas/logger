// +build !js

package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/shiena/ansicolor"
)

var output = ansicolor.NewAnsiColorWriter(os.Stdout)

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

		stackLine := fmt.Sprintf("\t%s:%d\n", file, line)
		stackFunc := fmt.Sprintf("\t\t%s\n", fnc.Name())
		if colored {
			stackLine = color.MagentaString(stackLine)
			stackFunc = color.MagentaString(stackFunc)
		}

		ret += stackLine + stackFunc
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
	printLog(color.FgRed, ErrPrefix, true, message, v...)
	os.Exit(1)
}

// Panic prints error log and call panic
func Panic(message string, v ...interface{}) {
	DisableBTrace()
	printLog(color.FgRed, ErrPrefix, true, message, v...)
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
		printLog(color.FgRed, ErrPrefix, true, "%v", err)
		return true
	}
	return false
}

// WarningErr prints err as warning log and returns true if err!=nil
func WarningErr(err error) bool {
	if err != nil {
		printLog(color.FgYellow, WarnPrefix, true, "%v", err)
		return true
	}
	return false
}

// FatalErr prints err as fatal log
func FatalErr(err error) {
	if err != nil {
		printLog(color.FgRed, ErrPrefix, true, "%v", err)
		os.Exit(1)
	}
}

// JSONDebug prints object in json format
func JSONDebug(a interface{}) {
	buf, err := json.MarshalIndent(a, "", "  ")
	if !WarningErr(err) {
		Debug("%s", buf)
	}
}
