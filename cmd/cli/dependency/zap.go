package dependency

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"xrf197ilz35aq0"
)

// https://github.com/uber-go/zap
// https://betterstack.com/community/guides/logging/go/zap/

type ZapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func (z ZapLogger) Info(message string) {
	z.logger.Info(message)
}

func (z ZapLogger) Warn(message string) {
	z.logger.Warn(message)
}

func (z ZapLogger) Debug(message string) {
	z.logger.Debug(message)
}

func (z ZapLogger) Error(message string) {
	z.logger.Error(message)
}

func (z ZapLogger) Fatal(message string) {
	z.logger.Fatal(message)
}

func (z ZapLogger) Panic(message string) {
	z.logger.Panic(message)
}

func NewZap(level string, dev bool, initialFields map[string]interface{}) (ZapLogger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(loggerLevel(level)),
		Development:       dev,
		Sampling:          nil,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: initialFields,
	}

	logger, err := config.Build()
	if err != nil {
		return ZapLogger{}, err
	}

	return ZapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}, nil
}

func CustomZapLogger(devEnv bool, level string, out io.Writer, initialFields []zapcore.Field) ZapLogger {
	// log outputs
	// log to multiple out puts. e.g file & console (os.Stdout)
	file := zapcore.AddSync(out)
	stdout := zapcore.AddSync(os.Stdout)

	zapLevel := loggerLevel(level)

	logLvl := zap.NewAtomicLevelAt(zapLevel)
	encoderConfig := createEncoderConfig(devEnv)

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	fileCore := zapcore.NewCore(fileEncoder, file, logLvl)
	consoleCore := zapcore.NewCore(consoleEncoder, stdout, logLvl)
	fileCore.With(initialFields)
	consoleCore.With(initialFields)
	// The NewTee() method duplicates log entries into two or more destinations.
	// In this case, the logs are sent to the standard output using a colorized plaintext format,
	// while the JSON equivalent is sent to the file
	core := zapcore.NewTee(fileCore, consoleCore)

	zapLogger := zap.New(core)
	zapLogger = zapLogger.WithOptions(zap.Fields(initialFields...))
	return ZapLogger{
		logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}
}

func createEncoderConfig(dev bool) zapcore.EncoderConfig {
	if dev {
		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		developmentCfg.TimeKey = "timestamp"
		return developmentCfg
	}
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	return productionCfg
}

func loggerLevel(level string) zapcore.Level {
	switch level {
	case xrf197ilz35aq0.DEBUG:
		return zapcore.DebugLevel
	case xrf197ilz35aq0.INFO:
		return zapcore.InfoLevel
	case xrf197ilz35aq0.WARN:
		return zapcore.WarnLevel
	case xrf197ilz35aq0.ERROR:
		return zapcore.ErrorLevel
	case xrf197ilz35aq0.FATAL:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
