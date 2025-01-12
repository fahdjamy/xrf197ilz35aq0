package internal

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
)

type Logger interface {
	Info(message string)
	Warn(message string)
	Debug(message string)
	Error(message string)
	Fatal(message string)
	SetPrefix(prefix string)
}
