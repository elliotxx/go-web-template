package safeutil

import (
	"os"
	"testing"
	"time"

	"github.com/elliotxx/safe"
	"github.com/sirupsen/logrus"
)

func TestGo(t *testing.T) {
	type args struct {
		do safe.DoFunc
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful-recover-crash-in-safe-Go",
			args: args{
				do: func() {
					panic("ah, I'm down")
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Go(tt.args.do)
			time.Sleep(time.Second * 1)
		})
	}
}

func TestGoL(t *testing.T) {
	type args struct {
		do     safe.DoFunc
		logger logrus.FieldLogger
	}

	getTestingLogger := func() logrus.FieldLogger {
		logger := logrus.New()
		logger.SetOutput(os.Stdout)
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000000",
			PrettyPrint:     true,
		})

		return logger
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful-recover-crash-in-safe-Go",
			args: args{
				do: func() {
					panic("ah, I'm down")
				},
				logger: getTestingLogger(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GoL(tt.args.do, tt.args.logger)
			time.Sleep(time.Second * 1)
		})
	}
}
