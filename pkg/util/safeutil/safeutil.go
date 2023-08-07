package safeutil

import (
	"runtime/debug"

	"github.com/elliotxx/safe"
	"github.com/sirupsen/logrus"
)

// LoggerRecoverHandler returns a recover handler by the given logger.
//
// Example:
//
//	func() {
//	  defer safe.HandleCrash(LoggerRecoverHandler(logrus.New()))
//	  ...
//	}
func LoggerRecoverHandler(logger logrus.FieldLogger) safe.RecoverHandler {
	return func(r any) {
		msgFormat := "Recovered as [%v] from stack: %s"

		if logger != nil {
			logger.Errorf(msgFormat, r, debug.Stack())
		}
	}
}

// Go starts a recoverable goroutine with a new logger (logrus.New()).
//
// Example:
//
//	safeutil.Go(func(){...})
func Go(do safe.DoFunc) {
	safe.GoR(do, LoggerRecoverHandler(logrus.New()))
}

// GoL starts a recoverable goroutine with a given logger.
//
// Example:
//
//	safeutil.GoL(func(){...}, logger)
func GoL(do safe.DoFunc, logger logrus.FieldLogger) {
	safe.GoR(do, LoggerRecoverHandler(logger))
}
