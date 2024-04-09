package slog

import (
	"fmt"
	"go-portfolios-tracker/internal/config"
	"io"
	"log/slog"
	"os"
)

// var h slog.Handler

type Logger struct {
	*slog.Logger
}

func NewLogger(env string) *Logger {
	var level slog.Level

	switch env {
	case config.EnvLocal:
		level = slog.LevelDebug
	case config.EnvDev:
		level = slog.LevelDebug
	case config.EnvProd:
		level = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	}

	err := os.MkdirAll("logs", 0770)
	if err != nil {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	// defer allFile.Close()

	w := io.MultiWriter(os.Stderr, allFile)

	h := slog.NewTextHandler(w, handlerOptions)
	return &Logger{slog.New(h)}
	// slog.SetDefault(slog.New(h))
}

func (l *Logger) Infof(template string, args ...any) {
	l.Logger.Info(fmt.Sprintf(template, args))
}

func (l *Logger) Info(args ...any) {
	l.Logger.Info(fmt.Sprint(args))
}

func (l *Logger) Debugf(template string, args ...any) {
	l.Logger.Debug(fmt.Sprintf(template, args))
}

func (l *Logger) Debug(args ...any) {
	l.Logger.Debug(fmt.Sprint(args))
}

func (l *Logger) Errorf(template string, args ...any) {
	l.Logger.Error(fmt.Sprintf(template, args))
}

func (l *Logger) Error(args ...any) {
	l.Logger.Error(fmt.Sprint(args))
}

func (l *Logger) Fatalf(template string, args ...any) {
	l.Logger.Error(fmt.Sprintf(template, args))
	os.Exit(1)
}

func (l *Logger) Fatal(args ...any) {
	l.Logger.Error(fmt.Sprint(args))
	os.Exit(1)
}
