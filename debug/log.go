package debug

import (
	"time"

	"github.com/Sirupsen/logrus"
)

// Time measures duration of a function
//
// Usage: (notice the suffix of '()')
//
// func A() {
//         defer debug.Time(log, "some %s", "args")()
//         ....
// }
func Time(log logrus.FieldLogger, format string, args ...interface{}) func() {
	if logrus.StandardLogger().Level < logrus.DebugLevel {
		return func() {}
	}
	start := time.Now()
	log.WithField("state", "started").Debugf(format, args...)
	return func() {
		log.WithField("state", "finished").
			WithField("duration", time.Now().Sub(start)).
			Debugf(format, args...)
	}
}
