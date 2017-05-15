package log

import (
	"io/ioutil"
	"os"
	"fmt"
	"time"
)

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var stdout = os.Stdout
var stderr = os.Stderr
var thresholdLevel = TRACE

func Trace(format string, v ...interface{}) {
	Log(TRACE, format, v...)
}

func Debug(format string, v ...interface{}) {
	Log(DEBUG, format, v...)
}

func Info(format string, v ...interface{}) {
	Log(INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
	Log(WARN, format, v...)
}

func Error(format string, v ...interface{}) {
	Log(ERROR, format, v...)
}

func Fatal(format string, v ...interface{}) {
	Log(FATAL, format, v...)
}

func Log(level LogLevel, format string, v ...interface{}) {
	writer := ioutil.Discard
	if (level >= thresholdLevel) {
		if (level == FATAL) {
			writer = stderr
		} else {
			writer = stdout
		}
	}
	now := time.Now().Format("2006-01-02 15:04:05.000")
	result := []interface{}{now, level.toString()}
	result = append(result, v...)

	fmt.Fprintf(writer, "%s - %s - " + format + "\n", result...)
}

func Stdout(writer *os.File) {
	stdout = writer
}

func Stderr(writer *os.File) {
	fmt.Sprintln()
	stderr = writer
}

func Level(level LogLevel) {
	thresholdLevel = level
}

func (level *LogLevel) toString() string {
	switch *level {
	case TRACE:
		return "[TRACE]";
	case DEBUG:
		return "[DEBUG]";
	case INFO:
		return "[INFO] ";
	case WARN:
		return "[WARN] ";
	case ERROR:
		return "[ERROR]";
	case FATAL:
		return "[FATAL]";
	default:
		return "[UNKNONW]"
	}
}