package dependency

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"xrf197ilz35aq0"
)

// https://github.com/uber-go/zap
// https://betterstack.com/community/guides/logging/go/zap/

type ZapLogger struct {
	prefix string
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func (z ZapLogger) Info(message string) {
	z.logger.Info(message)
}

func (z ZapLogger) Warn(message string) {
	z.sugar.Warn(message)
}

func (z ZapLogger) Debug(message string) {
	z.sugar.Debug(message)
}

func (z ZapLogger) Error(message string) {
	z.sugar.Error(message)
}

func (z ZapLogger) Fatal(message string) {
	z.sugar.Fatal(message)
}

func (z ZapLogger) Panic(message string) {
	z.sugar.Panic(message)
}

func NewZap(level string, dev bool, _ io.Writer) (ZapLogger, error) {
	_, err := loggerLevel(level)
	if err != nil {
		return ZapLogger{}, err
	}

	logger, err := zap.NewProduction()

	if err != nil {
		return ZapLogger{}, err
	}

	if dev {
		logger = logger.WithOptions(zap.Development())
	}
	//logger = logger.WithOptions(zap.Fields(zap.Object(initialFields)))

	return ZapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}, nil
}

func loggerLevel(level string) (zapcore.Level, error) {
	switch level {
	case xrf197ilz35aq0.DEBUG:
		return zapcore.DebugLevel, nil
	case xrf197ilz35aq0.INFO:
		return zapcore.InfoLevel, nil
	case xrf197ilz35aq0.WARN:
		return zapcore.WarnLevel, nil
	case xrf197ilz35aq0.ERROR:
		return zapcore.ErrorLevel, nil
	case xrf197ilz35aq0.FATAL:
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}
