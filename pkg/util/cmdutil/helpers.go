package cmdutil

import (
	"os"
	"strings"

	"github.com/elliotxx/go-web-template/pkg/util/pretty"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

const (
	DefaultErrorExitCode = 1
)

var fatalErrHandler = fatal

// fatal prints the message (if provided) and then exits.
func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		pretty.ErrorT.WithWriter(os.Stderr).Print(msg)
	}
	os.Exit(code)
}

// ErrExit may be passed to CheckError to instruct it to output nothing but exit with
// status code 1.
var ErrExit = errors.Errorf("exit")

// CheckErr prints a user friendly error to STDERR and exits with a non-zero
// exit code. Unrecognized errors will be printed with an "error: " prefix.
//
// This method is generic to the command in use and may be used by non-Kubectl
// commands.
func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

// checkErr formats a given error as a string and calls the passed handleErr
// func with that string and an kubectl exit code.
func checkErr(err error, handleErr func(string, int)) {
	// flatten errors
	if merr, ok := err.(*multierror.Error); ok {
		err = multierror.Flatten(merr.ErrorOrNil())
	}

	if err == nil {
		return
	}

	switch {
	case err == ErrExit:
		handleErr("", DefaultErrorExitCode)
	default:
		handleErr(err.Error(), DefaultErrorExitCode)
	}
}
