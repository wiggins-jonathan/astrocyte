package cmd

import (
	"io"
	"log/slog"
	"os"
)

type logHandler struct {
	HandlerOptions slog.HandlerOptions
	Output         io.Writer
	Format         string
}

type LoggerOption func(*logHandler)

func NewLogger(options ...LoggerOption) *slog.Logger {
	h := &logHandler{
		Output:         os.Stdout,
		Format:         "text",
		HandlerOptions: slog.HandlerOptions{Level: slog.LevelInfo},
	}

	for _, option := range options {
		option(h)
	}

	var handler slog.Handler
	switch h.Format {
	case "json":
		handler = slog.NewJSONHandler(h.Output, &h.HandlerOptions)
	default:
		handler = slog.NewTextHandler(h.Output, &h.HandlerOptions)
	}

	return slog.New(handler)
}

func WithLogLevel(level string) LoggerOption {
	return func(h *logHandler) {
		switch level {
		case "debug":
			h.HandlerOptions.Level = slog.LevelDebug
		case "info":
			h.HandlerOptions.Level = slog.LevelInfo
		case "warn":
			h.HandlerOptions.Level = slog.LevelWarn
		case "error":
			h.HandlerOptions.Level = slog.LevelError
		default:
			h.HandlerOptions.Level = slog.LevelInfo
		}
	}
}

func WithOutput(w io.Writer) LoggerOption {
	return func(h *logHandler) {
		h.Output = w
	}
}

func WithFormat(format string) LoggerOption {
	return func(h *logHandler) {
		h.Format = format
	}
}
