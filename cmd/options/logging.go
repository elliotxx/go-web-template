package options

import (
	"io"
	"path/filepath"
	"runtime"
	"strconv"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/elliotxx/errors"
	"github.com/elliotxx/go-web-template/cmd/options/types"
	"github.com/elliotxx/go-web-template/pkg/server"
	"github.com/hashicorp/go-multierror"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ types.Options = &LoggingOptions{}

// LoggingOptions provides the logging configuration.
type LoggingOptions struct {
	LogLevel            string `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	DisableText         bool   `json:"disableText,omitempty" yaml:"disableText,omitempty"`
	TextPretty          bool   `json:"textPretty,omitempty" yaml:"textPretty,omitempty"`
	JSONPretty          bool   `json:"jsonPretty,omitempty" yaml:"jsonPretty,omitempty"`
	ReportCaller        bool   `json:"reportCaller,omitempty" yaml:"reportCaller,omitempty"`
	DumpCurrentConfig   bool   `json:"dumpCurrentConfig,omitempty" yaml:"dumpCurrentConfig,omitempty"`
	EnableLoggingToFile bool   `json:"enableLoggingToFile,omitempty" yaml:"enableLoggingToFile,omitempty"`
	LoggingDirectory    string `json:"loggingDirectory,omitempty" yaml:"loggingDirectory,omitempty"`
}

// NewLoggingOptions returns a LoggingOptions instance with the default values
func NewLoggingOptions() *LoggingOptions {
	return &LoggingOptions{
		DumpCurrentConfig: true,
		TextPretty:        true,
		LogLevel:          "info",
	}
}

// ApplyTo apply logging options to the server config
func (o *LoggingOptions) ApplyTo(config *server.Config) {
	config.LoggingDirectory = o.LoggingDirectory
}

// Validate checks LoggingOptions and return a slice of found error(s)
func (o *LoggingOptions) Validate() error {
	if o == nil {
		return errors.Errorf("options is nil")
	}

	var err *multierror.Error

	if _, err2 := logrus.ParseLevel(o.LogLevel); err2 != nil {
		err = multierror.Append(err, err2)
	}

	if o.JSONPretty && !o.DisableText {
		err = multierror.Append(err, errors.Errorf("--json-pretty cannot be enabled when --disable-text is disabled"))
	}

	if o.EnableLoggingToFile && o.LoggingDirectory == "" {
		err = multierror.Append(err, errors.Errorf("--logging-directory must be not empty when --enable-logging-to-file is enabled"))
	}

	if o.EnableLoggingToFile && o.TextPretty {
		err = multierror.Append(err, errors.Errorf("--text-pretty cannot be enabled when --enable-logging-to-file is enabled"))
	}

	return err.ErrorOrNil()
}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *LoggingOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVarP(&o.LogLevel, "log-level", "L", o.LogLevel, "Log level. Valid values: [trace, debug, info, warn, warning, error, fatal, panic]")
	fs.BoolVar(&o.JSONPretty, "json-pretty", o.JSONPretty, "JSONPretty will indent all json logs")
	fs.BoolVar(&o.DisableText, "disable-text", o.DisableText, "Disable text mode")
	fs.BoolVar(&o.TextPretty, "text-pretty", o.TextPretty, "TextPretty will colorize the logs")
	fs.BoolVar(&o.ReportCaller, "report-caller", o.ReportCaller, "Flag for whether to log caller info")
	fs.BoolVar(&o.DumpCurrentConfig, "dump-current-config", o.DumpCurrentConfig, "Dump current configuration")
	fs.BoolVar(&o.EnableLoggingToFile, "enable-logging-to-file", o.EnableLoggingToFile, "Enable logging to file")
	fs.StringVar(&o.LoggingDirectory, "logging-directory", ".", "Specify which directory to log to")
}

// helper function configures the logging.
func (o *LoggingOptions) InitLogging(projectName string) error {
	if o == nil {
		return errors.Errorf("logging options is nil")
	}

	if len(o.LogLevel) != 0 {
		logLevel, err := logrus.ParseLevel(o.LogLevel)
		if err != nil {
			return errors.Wrap(err, "failed to init logging configuration")
		}

		logrus.SetLevel(logLevel)
	}

	if o.ReportCaller {
		logrus.SetReportCaller(true)
	}

	var formatter logrus.Formatter
	if o.DisableText {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000000",
			PrettyPrint:     o.JSONPretty,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				return "", filepath.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			},
		}
	} else {
		formatter = &nested.Formatter{
			TimestampFormat: "2006-01-02 15:04:05.000000",
			NoColors:        !o.TextPretty,
			CallerFirst:     true,
			CustomCallerFormatter: func(frame *runtime.Frame) string {
				return " " + filepath.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			},
		}
	}
	logrus.SetFormatter(formatter)

	// Set log output to file
	if o.EnableLoggingToFile {
		// Set output logs to files (error, info, trace)
		setLoggingToFile(o.LoggingDirectory, formatter, projectName)
	}

	logrus.Info("Successfully initialized log configuration")

	return nil
}

// setLoggingToFile sets the logging entry to the specified file(s)
func setLoggingToFile(loggingDirectory string, writerFormatter logrus.Formatter, projectName string) {
	// Create lumberjack writers for three different log levels
	errorLoggingFile := filepath.Join(loggingDirectory, projectName+".error.log")
	errorRotateWriter := &lumberjack.Logger{
		Filename:   errorLoggingFile,
		MaxSize:    1,  // MaxSize the max size in MB of the logfile before it's rolled
		MaxBackups: 10, // MaxBackups the max number of rolled files to keep
		MaxAge:     30, // MaxAge the max age in days to keep a logfile
	}
	infoLoggingFile := filepath.Join(loggingDirectory, projectName+".log")
	infoRotateWriter := &lumberjack.Logger{
		Filename:   infoLoggingFile,
		MaxSize:    1,
		MaxBackups: 10,
		MaxAge:     30,
	}
	traceLoggingFile := filepath.Join(loggingDirectory, projectName+".trace.log")
	traceRotateWriter := &lumberjack.Logger{
		Filename:   traceLoggingFile,
		MaxSize:    1,
		MaxBackups: 10,
		MaxAge:     30,
	}

	// Set output logs to files (error, info, trace)
	writerMap := lfshook.WriterMap{
		logrus.PanicLevel: io.MultiWriter(errorRotateWriter, infoRotateWriter, traceRotateWriter),
		logrus.FatalLevel: io.MultiWriter(errorRotateWriter, infoRotateWriter, traceRotateWriter),
		logrus.ErrorLevel: io.MultiWriter(errorRotateWriter, infoRotateWriter, traceRotateWriter),
		logrus.WarnLevel:  io.MultiWriter(infoRotateWriter, traceRotateWriter),
		logrus.InfoLevel:  io.MultiWriter(infoRotateWriter, traceRotateWriter),
		logrus.DebugLevel: io.MultiWriter(traceRotateWriter),
		logrus.TraceLevel: io.MultiWriter(traceRotateWriter),
	}

	// Finally, log output to stderr, error file, info file, and trace file
	logrus.AddHook(lfshook.NewHook(
		writerMap,
		writerFormatter,
	))

	logrus.WithFields(logrus.Fields{
		"logging-directory":  loggingDirectory,
		"project-name":       projectName,
		"error-logging-file": errorLoggingFile,
		"info-logging-file":  infoLoggingFile,
		"trace-logging-file": traceLoggingFile,
	}).Debug("Successfully set logging to file")
}
