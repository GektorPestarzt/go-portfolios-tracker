package logging

import (
	"fmt"
	"go-portfolios-tracker/internal/config"
	"io"
	"log/slog"
	"os"
)

var h slog.Handler

type Logger struct {
	*slog.Logger
}

func GetLogger() *Logger {
	return &Logger{slog.New(h)}
}

func MustLoad(env string) {
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
		AddSource: true,
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

	h = slog.NewTextHandler(w, handlerOptions)
	slog.SetDefault(slog.New(h))
}

func (l *Logger) Fatal(args ...any) {
	l.Error(fmt.Sprint(args))
	os.Exit(1)
}
