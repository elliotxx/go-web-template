package ctxutil

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ContextKey string

const (
	ContextKeyLogger ContextKey = "logger"
)

// GetLogger returns the logger from the given context.
//
// Example:
//
//	logger := ctxutil.GetLogger(ctx)
func GetLogger(ctx context.Context) logrus.FieldLogger {
	if logger, ok := ctx.Value(ContextKeyLogger).(logrus.FieldLogger); ok {
		return logger
	}

	return logrus.New()
}

// WithLogger returns a context by the TODO context and the given logger.
//
// Example:
//
//	ctx = ctxutil.WithLogger(logger)
func WithLogger(logger logrus.FieldLogger) context.Context {
	return context.WithValue(context.TODO(), ContextKeyLogger, logger)
}

// CtxWithLogger returns a context by the parent context and the given logger.
//
// Example:
//
//	ctx = ctxutil.CtxWithLogger(ctx, logger)
func CtxWithLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, logger)
}
