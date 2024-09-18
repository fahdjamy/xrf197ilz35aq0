package main

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strconv"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/cmd"
	"xrf197ilz35aq0/dependency"
	"xrf197ilz35aq0/internal/random"
)

func main() {
	environment := cmd.GetEnvironment()

	health := xrf197ilz35aq0.NewHealth()
	logFileOutPut := &lumberjack.Logger{
		Filename:   ".logs/xrf197ilz.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	}

	requestId, err := cmd.GenerateRequestId()
	if err != nil {
		requestId = strconv.Itoa(int(random.PositiveInt64()))
	}
	initialFields := []zap.Field{
		zap.String("requestId", requestId),
		zap.String("version", health.Version()),
	}

	config, err := xrf197ilz35aq0.NewConfig(environment.Name)
	if err != nil {
		panic(err)
	}

	dbUri := mongoUri(config)
	fmt.Println(dbUri)
	logger := dependency.CustomZapLogger(environment.LogMode, config.Log.Level, logFileOutPut, initialFields)
	logger.Info(fmt.Sprintf("application version '%s'", health.Version()))
}

func mongoUri(config xrf197ilz35aq0.Config) string {
	mongoConfig := config.Database.Mongo
	baseUri := os.Getenv(mongoConfig.Uri)
	return fmt.Sprintf("%sretryWrites=%tw=%sappName=%s",
		baseUri,
		mongoConfig.RetryWrites,
		mongoConfig.Acknowledgment,
		mongoConfig.AppName,
	)
}
