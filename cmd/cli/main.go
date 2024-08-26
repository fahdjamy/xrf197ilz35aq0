package main

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"strconv"
	"time"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/cmd/cli/dependency"
	"xrf197ilz35aq0/random"
)

func main() {
	health := xrf197ilz35aq0.NewHealth()
	logFileOutPut := &lumberjack.Logger{
		Filename:   ".logs/xrf197ilz.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	}

	requestId, err := generateRequestId()
	if err != nil {
		requestId = strconv.Itoa(int(random.PositiveInt64()))
	}
	initialFields := []zap.Field{
		zap.String("requestId", requestId),
		zap.String("version", health.Version()),
	}

	logger := dependency.CustomZapLogger(true, xrf197ilz35aq0.DEBUG, logFileOutPut, initialFields)
	logger.Info(fmt.Sprintf("application version '%s'", health.Version()))
}

func generateRequestId() (string, error) {
	uniqueStr, err := random.TimeBasedString(time.Now().Unix(), 21)
	if err != nil {
		return "", err
	}

	uniqueInt64 := random.PositiveInt64()
	uniqueInt64Str := strconv.Itoa(int(uniqueInt64))

	if len(uniqueInt64Str) > 10 {
		uniqueInt64Str = uniqueInt64Str[2:]
	}

	partStr := uniqueStr[0:12]

	return fmt.Sprintf("%s.%s", uniqueInt64Str, partStr), nil
}
