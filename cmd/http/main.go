package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/cmd"
	"xrf197ilz35aq0/core/service/user"
	"xrf197ilz35aq0/dependency"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/server/http"
	"xrf197ilz35aq0/storage/mongo"
)

func main() {
	// get the globally set environment variables
	environment := cmd.GetEnvironment()

	// get the configuration for the application
	config, err := xrf.NewConfig(environment.Name)
	if err != nil {
		panic(err)
	}

	// get the health information about the application
	health := xrf.NewHealth()

	// create a logger
	logFileOutPut := &lumberjack.Logger{
		Filename:   config.Log.Filename,
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	}
	initialFields := []zap.Field{
		zap.String("os", health.Runtime.OS),
	}
	logPrefix := fmt.Sprintf("requestId='%s'", cmd.GenerateRequestId())
	logger := dependency.CustomZapLogger(environment.LogMode, config.Log.Level, logFileOutPut, logPrefix, initialFields)
	logger.Info(fmt.Sprintf("appVersion='%s' :: os='%s' :: message='application starting...'", health.Version(), health.Runtime.OS))

	// connect to the Mongo Database
	dbConnStr, err := mongoUri(config)
	backgroundCtx := context.Background()
	if err != nil {
		logger.Panic(fmt.Sprintf("appStarted=false :: err%s", err.Error()))
		return
	}
	databaseName := config.Database.Mongo.DatabaseName
	mongoClient, err := mongo.NewClient(backgroundCtx, dbConnStr, databaseName)
	if err != nil {
		internalError := xrfErr.Internal{Err: err, Message: "failed to connect to mongo"}
		logger.Panic(fmt.Sprintf("appStarted=false :: %s", internalError.Error()))
	}
	mongoDB := mongo.NewStore(logger, mongoClient, databaseName, backgroundCtx)

	// create services
	settingsManager := user.NewSettingManager(logger)
	userManager := user.NewUserManager(logger, settingsManager, mongoDB)

	// create the router and start the server
	router := mux.NewRouter().StrictSlash(true)
	apiServer := http.NewHttpServer(logger, router, config, userManager, backgroundCtx)
	apiServer.Start()
}

func mongoUri(config xrf.Config) (string, error) {
	mongoConfig := config.Database.Mongo
	baseUri := os.Getenv(mongoConfig.Uri)
	if baseUri == "" {
		return "", &xrfErr.Internal{
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
