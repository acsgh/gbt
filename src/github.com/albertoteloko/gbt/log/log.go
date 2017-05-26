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

func Tracef(format string, v ...interface{}) {
	Logf(TRACE, format, v...)
}

func Debugf(format string, v ...interface{}) {
	Logf(DEBUG, format, v...)
}

func Infof(format string, v ...interface{}) {
	Logf(INFO, format, v...)
}

func Warnf(format string, v ...interface{}) {
	Logf(WARN, format, v...)
}

func Errorf(format string, v ...interface{}) {
	Logf(ERROR, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	Logf(FATAL, format, v...)
}

func Logf(level LogLevel, format string, v ...interface{}) {
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