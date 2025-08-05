package log

import (
	"log/slog"
)

type Logger struct {
	*slog.Logger
}


func (l *Logger) WithModule(module string) *Logger {
	return &Logger{l.With(slog.String("module", module))}
}
