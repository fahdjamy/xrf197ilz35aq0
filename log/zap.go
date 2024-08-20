package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

// https://github.com/uber-go/zap

type zapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func (z zapLogger) Info(args ...any) {
	z.sugar.Info(args)
}

func (z zapLogger) Warn(args ...interface{}) {
	z.sugar.Warn(args...)
}

func (z zapLogger) Debug(args ...interface{}) {
	z.sugar.Debug(args...)
}

func (z zapLogger) Error(args ...interface{}) {
	z.sugar.Error(args...)
}

func (z zapLogger) Fatal(args ...interface{}) {
	z.sugar.Fatal(args...)
}

func (z zapLogger) Panic(args ...interface{}) {
	z.sugar.Panic(args...)
}

func (z zapLogger) InfoF(format string, args ...any) {
	z.sugar.Infof(format, args...)
}

func (z zapLogger) WarnF(format string, args ...any) {
	z.sugar.Warnf(format, args...)
}

func (z zapLogger) DebugF(format string, args ...any) {
	z.sugar.Debugf(format, args...)
}

func (z zapLogger) Errorf(format string, args ...any) {
	z.sugar.Errorf(format, args...)
}

func (z zapLogger) Fatalf(format string, args ...any) {
	z.sugar.Fatalf(format, args...)
}

func (z zapLogger) PanicF(format string, args ...any) {
	z.sugar.Panicf(format, args...)
}

func newZap(level string, dev bool, _ io.Writer) (zapLogger, error) {
	_, err := loggerLevel(level)
	if err != nil {
		return zapLogger{}, err
	}

	logger, err := zap.NewProduction()

	if err != nil {
		return zapLogger{}, err
	}

	if dev {
		logger = logger.WithOptions(zap.Development())
	}

	return zapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}, nil
}

func loggerLevel(level string) (zapcore.Level, error) {
	switch level {
	case DEBUG:
		return zapcore.DebugLevel, nil
	case INFO:
		return zapcore.InfoLevel, nil
	case WARN:
		return zapcore.WarnLevel, nil
	case ERROR:
		return zapcore.ErrorLevel, nil
	case FATAL:
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}
