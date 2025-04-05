package cmd

import (
	"io"
	"log/slog"
	"os"
)

type logHandler struct {
	handlerOptions slog.HandlerOptions
	output         io.Writer
	format         string
}

type LoggerOption func(*logHandler)

func NewLogger(options ...LoggerOption) *slog.Logger {
	h := &logHandler{
		output:         os.Stdout,
		format:         "text",
		handlerOptions: slog.HandlerOptions{Level: slog.LevelInfo},
	}

	for _, option := range options {
		option(h)
	}

	var handler slog.Handler
	switch h.format {
	case "json":
		handler = slog.NewJSONHandler(h.output, &h.handlerOptions)
	default:
		handler = slog.NewTextHandler(h.output, &h.handlerOptions)
	}

	return slog.New(handler)
}

func WithLogLevel(level string) LoggerOption {
	return func(h *logHandler) {
		switch level {
		case "debug":
			h.handlerOptions.Level = slog.LevelDebug
		case "info":
			h.handlerOptions.Level = slog.LevelInfo
		case "warn":
			h.handlerOptions.Level = slog.LevelWarn
		case "error":
			h.handlerOptions.Level = slog.LevelError
		default:
			h.handlerOptions.Level = slog.LevelInfo
		}
	}
}

func WithOutput(w io.Writer) LoggerOption {
	return func(h *logHandler) {
		h.output = w
	}
}

func WithFormat(format string) LoggerOption {
	return func(h *logHandler) {
		h.format = format
	}
}
