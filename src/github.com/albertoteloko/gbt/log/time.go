package log

import (
	"time"
)

func LogTime(taskName string, task func()) {
	startTime := time.Now().UnixNano()
	Debug("%s start", taskName)
	task()
	Info("%s in %d ms", taskName, (time.Now().UnixNano()-startTime)/1000000)
}
