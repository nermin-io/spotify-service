package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/nermin-io/spotify-service/apiserver"
	"github.com/nermin-io/spotify-service/spotify"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var debug bool

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.Parse()

	logger, err := initLogger()
	if err != nil {
		return err
	}
	defer logger.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	spotifyClient := spotify.NewClient(logger)
	api := apiserver.NewHandler(logger, spotifyClient)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: api,
	}

	go func() {
		logger.Sugar().Infof("listening on port %s", port)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server crashed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}

func initLogger() (*zap.Logger, error) {
	logLevel := zap.InfoLevel
	if debug {
		logLevel = zap.DebugLevel
	}
	return zap.Config{
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
}