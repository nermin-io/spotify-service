package logging

import (
	"context"
	"github.com/nermin-io/spotify-service/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

const traceLogField = "logging.googleapis.com/trace"

func Init(debug bool) (*zap.Logger, error) {
	logLevel := zap.InfoLevel
	if debug {
		logLevel = zap.DebugLevel
	}
	logger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "severity",
			MessageKey:     "message",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	globalLogger = logger
	return logger, err
}

// FromContext accepts a context.Context and returns a child Logger that provides additional fields
// such as the Google Cloud Logging trace.
func FromContext(ctx context.Context) *zap.Logger {
	traceID := trace.FromContext(ctx)
	if traceID == "" {
		return globalLogger
	}

	return globalLogger.With(zap.String(traceLogField, traceID))
}
