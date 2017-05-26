package log

import (
	"time"
)

func LogTime(taskName string, task func()) {
	startTime := time.Now().UnixNano()
	Debugf("%s start", taskName)
	task()
	Infof("%s in %d ms", taskName, (time.Now().UnixNano()-startTime)/1000000)
}

func LogTimeWithError(taskName string, task func() error) error {
	startTime := time.Now().UnixNano()
	Debugf("%s start", taskName)
	err := task()
	Infof("%s in %d ms", taskName, (time.Now().UnixNano()-startTime)/1000000)
	return err
}
