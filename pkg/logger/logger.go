package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(ho *slog.HandlerOptions) *Logger {
	return &Logger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, ho)),
	}
}

func (l *Logger) Info(msg string, args ...any) {
	fmt.Print("\033[32m")
	l.logger.Info(msg, args...)
	fmt.Print("\033[0m")
}

func (l *Logger) Error(msg string, args ...any) {
	fmt.Print("\033[91m")
	l.logger.Info(msg, args...)
	fmt.Print("\033[0m")
}

func (l *Logger) Warn(msg string, args ...any) {
	fmt.Print("\033[93m")
	l.logger.Info(msg, args...)
	fmt.Print("\033[0m")
}

func (l *Logger) Debug(msg string, args ...any) {
	fmt.Print("\033[95m")
	l.logger.Info(msg, args...)
	fmt.Print("\033[0m")
}
