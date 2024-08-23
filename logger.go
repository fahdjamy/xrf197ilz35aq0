package xrf197ilz35aq0

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
	FATAL = "fatal"
)

type Logger interface {
	Info(message string)
	Warn(message string)
	Debug(message string)
	Error(message string)
	Fatal(message string)
	Panic(message string)
}
