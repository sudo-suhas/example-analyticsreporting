package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// TimeTrack , surprise surprise, tracks time for any function's execution
// Can be used with a simple defer:
//     defer TimeTrack(time.Now(), "Task Name")
func TimeTrack(start time.Time, name string) {
	log.WithField("ExecTime", time.Since(start)).Debugf("Time track for '%s'", name)
}
