package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/cmd"
	"xrf197ilz35aq0/cmd/cli/adapter"
	"xrf197ilz35aq0/dependency"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/storage/mongo"
)

func main() {
	// get the globally set environment variables
	environment := cmd.GetEnvironment()

	// get the configuration for the application
	config, err := xrf197ilz35aq0.NewConfig(environment.Name)
	if err != nil {
		panic(err)
	}

	// get the health information about the application
	health := xrf197ilz35aq0.NewHealth()

	// create a logger
	logFileOutPut := &lumberjack.Logger{
		Filename:   ".logs/xrf197ilz.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	}

	initialFields := []zap.Field{
		zap.String("os", health.Runtime.OS),
	}

	logPrefix := fmt.Sprintf("appVersion='%s' :: requestId='%s'", health.Version(), cmd.GenerateRequestId())
	logger := dependency.CustomZapLogger(environment.LogMode, config.Log.Level, logFileOutPut, logPrefix, initialFields)

	// connect to the Mongo Database
	dbConnStr, err := mongoUri(config)
	if err != nil {
		logger.Panic(fmt.Sprintf("appStarted=false :: err%s", err.Error()))
		return
	}
	databaseName := config.Database.Mongo.DatabaseName
	mongoClient, err := mongo.NewClient(context.Background(), dbConnStr, databaseName)

	if err != nil {
		internalError := xrfErr.Internal{
			Err:     err,
			Time:    time.Now(),
			Source:  "cmd/cli/main",
			Message: "failed to connect to mongo",
		}
		logger.Panic(fmt.Sprintf("appStarted=false :: err%s", internalError.Error()))
		return
	}

	// create a mongo store
	mongoStore := mongo.NewStore(logger, mongoClient, databaseName, context.Background())

	// create the app
	app := adapter.NewApp(logger, mongoStore)

	// start application
	app.Start()
}

func mongoUri(config xrf197ilz35aq0.Config) (string, error) {
	mongoConfig := config.Database.Mongo
	baseUri := os.Getenv(mongoConfig.Uri)
	if baseUri == "" {
		return "", &xrfErr.Internal{
			Time:    time.Now(),
			Source:  "cmd/cli/main#mongoUri",
			Message: "missing environment variable $" + mongoConfig.Uri,
		}
	}
	return fmt.Sprintf("%s?directConnection=%t&retryWrites=%t&w=%s&appName=%s",
		baseUri,
		mongoConfig.DirectConnection,
		mongoConfig.RetryWrites,
		mongoConfig.Acknowledgment,
		mongoConfig.AppName,
	), nil
}
