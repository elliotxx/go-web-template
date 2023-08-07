package ctxutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGetLogger(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	isWantedLogger := func(got any) bool {
		return got != nil
	}
	tests := []struct {
		name     string
		args     args
		wantFunc func(want any) bool
	}{
		{
			name: "successful-get-logger-from-context",
			args: args{
				ctx: context.WithValue(context.TODO(), ContextKeyLogger, logrus.WithError(fmt.Errorf("for test"))),
			},
			wantFunc: isWantedLogger,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLogger(tt.args.ctx); !isWantedLogger(got) {
				t.Errorf("GetLogger() = %v, is not the desired value", got)
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	type args struct {
		logger logrus.FieldLogger
	}
	isWantedContext := func(got context.Context) bool {
		if _, ok := got.Value(ContextKeyLogger).(logrus.FieldLogger); ok {
			return true
		}
		return false
	}
	tests := []struct {
		name     string
		args     args
		wantFunc func(got context.Context) bool
	}{
		{
			name: "successful-create-new-context-using-logger",
			args: args{
				logger: logrus.WithError(fmt.Errorf("for test")),
			},
			wantFunc: isWantedContext,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithLogger(tt.args.logger); !isWantedContext(got) {
				t.Errorf("GetLogger() = %v, is not the desired value", got)
			}
		})
	}
}
