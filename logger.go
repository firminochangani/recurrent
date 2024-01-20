package schedule

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, arg ...any)
	Infof(msg string, arg ...any)
}

type DefaultLogger struct {
	logger *slog.Logger
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (d *DefaultLogger) Info(msg string, arg ...any) {
	d.logger.Info(msg, arg...)
}

func (d *DefaultLogger) Infof(msg string, arg ...any) {
	d.logger.Info(fmt.Sprintf(msg, arg...))
}
