package xrf197ilz35aq0

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
	Panic(message string)
}
