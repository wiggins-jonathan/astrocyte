package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

type logHandler struct {
	handlerOptions slog.HandlerOptions
	output         io.Writer
	format         string
}

type LoggerOption func(*logHandler) error

func NewLogger(options ...LoggerOption) (*slog.Logger, error) {
	h := &logHandler{
		output:         os.Stdout,
		format:         "text",
		handlerOptions: slog.HandlerOptions{Level: slog.LevelInfo},
	}

	for _, option := range options {
		if err := option(h); err != nil {
			return nil, err
		}
	}

	var handler slog.Handler
	switch h.format {
	case "json":
		handler = slog.NewJSONHandler(h.output, &h.handlerOptions)
	default:
		handler = slog.NewTextHandler(h.output, &h.handlerOptions)
	}

	return slog.New(handler), nil
}

func WithLevel(level string) LoggerOption {
	return func(h *logHandler) error {
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
			return fmt.Errorf("%s is not a valid log level", level)
		}

		return nil
	}
}

func WithOutput(w io.Writer) LoggerOption {
	return func(h *logHandler) error {
		h.output = w
		return nil
	}
}

func WithFormat(format string) LoggerOption {
	return func(h *logHandler) error {
		switch format {
		case "json":
			h.format = "json"
		case "text":
			h.format = "text"
		default:
			return fmt.Errorf("%s is not a valid log format", format)
		}

		return nil
	}
}
