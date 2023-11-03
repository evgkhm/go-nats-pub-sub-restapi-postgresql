package logging

import (
	"log/slog"
	"os"
)

const (
	defaultLevel     = slog.LevelInfo
	defaultIsJSON    = false
	defaultAddSource = false
)

func Logger(opts ...LoggerOption) *slog.Logger {
	config := &LoggerOptions{
		Level:     defaultLevel,
		IsJSON:    defaultIsJSON,
		AddSource: defaultAddSource,
	}

	for _, opt := range opts {
		opt(config)
	}

	if err := config.Validate(); err != nil {
		return slog.Default()
	}

	options := &slog.HandlerOptions{
		AddSource: config.AddSource,
		Level:     config.Level,
	}

	var h slog.Handler = slog.NewTextHandler(os.Stdout, options)

	if config.IsJSON {
		h = slog.NewJSONHandler(os.Stdout, options)
	}

	return slog.New(h)
}

type LoggerOptions struct {
	Level     slog.Level
	AddSource bool
	IsJSON    bool
}

func (o *LoggerOptions) Validate() error {
	return nil
}

type LoggerOption func(*LoggerOptions)

// WithLevel logger option sets the log level, if not set, the default level is Info.
func WithLevel(level string) LoggerOption {
	return func(o *LoggerOptions) {
		var l slog.Level
		if err := l.UnmarshalText([]byte(level)); err != nil {
			l = slog.LevelInfo
		}

		o.Level = l
	}
}

func WithAddSource(addSource bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.AddSource = addSource
	}
}

func WithIsJSON(isJSON bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.IsJSON = isJSON
	}
}

//func L(ctx context.Context) *slog.Logger {
//	return loggerFromContext(ctx)
//}
