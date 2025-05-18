package logging

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

func Init(environment string) {
	var handler slog.Handler

	if environment == "development" {
		opts := &tint.Options{Level: slog.LevelDebug}
		if runtime.GOOS == "windows" {
			handler = tint.NewHandler(colorable.NewColorable(os.Stderr), opts)
		} else {
			handler = tint.NewHandler(os.Stderr, opts)
		}
	} else {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
