package log

import (
	"io"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
	FATAL = "fatal"
)

type Logger interface {
	Info(args ...any)
	InfoF(format string, args ...any)
	Warn(args ...interface{})
	WarnF(format string, args ...any)
	Debug(args ...interface{})
	DebugF(format string, args ...any)
	Error(args ...interface{})
	Errorf(format string, args ...any)
	Fatal(args ...interface{})
	Fatalf(format string, args ...any)
	Panic(args ...interface{})
	PanicF(format string, args ...any)
}

func New(writer io.Writer, level string) (Logger, error) {
	dev := false

	if level == DEBUG {
		dev = true
	}
	logger, err := newZap(level, dev, writer)

	return logger, err
}
