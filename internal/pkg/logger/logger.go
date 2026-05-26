package logger

import "log/slog"

const (
	keyError     = "error"
	keyComponent = "component"
)

func WithComponent(name string) *slog.Logger {
	return slog.Default().With(slog.String(keyComponent, name))
}

func Error(err error) slog.Attr {
	return slog.String(keyError, err.Error())
}
