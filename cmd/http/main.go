package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/core/repository"
	"xrf197ilz35aq0/core/service"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/dependency"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/server/http"
	"xrf197ilz35aq0/storage/mongo"
)

func main() {
	// get the globally set environment variables
	environment := internal.GetEnvironment()

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
	logPrefix := fmt.Sprintf("requestId='%s'", internal.GenerateRequestId())
	logger := dependency.CustomZapLogger(environment.LogMode, config.Log.Level, logFileOutPut, logPrefix, initialFields)
	logger.Info(fmt.Sprintf("appVersion='%s' :: os='%s' :: message='application starting...'", health.Version(), health.Runtime.OS))

	// connect to the Mongo Database
	dbConnStr, err := mongoUri(config)
	backgroundCtx := context.Background()
	if err != nil {
		logger.Error(fmt.Sprintf("appStarted=false :: err%s", err.Error()))
		return
	}
	databaseName := config.Database.Mongo.DatabaseName
	mongoClient, err := mongo.NewClient(backgroundCtx, dbConnStr, databaseName)
	if err != nil {
		internalError := xrfErr.Internal{Err: err, Message: "failed to connect to mongo"}
		logger.Error(fmt.Sprintf("appStarted=failure :: %s", internalError.Error()))
		return
	}
	// connect to MongoDB
	mongoDB := mongo.NewStore(logger, mongoClient, databaseName, backgroundCtx)
	logger.Debug(fmt.Sprintf("message='successfully connected to MongoDB' :: dbName=%s", databaseName))

	// create repositories
	userRepo := repository.NewUserRepository(mongoDB)
	settingRepo := repository.NewSettingsRepository(mongoDB)

	// create services
	settingsService := service.NewSettingService(logger, mongoDB, settingRepo, backgroundCtx)
	userService := service.NewUserService(logger, settingsService, mongoDB, userRepo, backgroundCtx)

	// create the router and start the server
	router := mux.NewRouter().StrictSlash(true)
	server := http.NewHttpServer(logger, router, config, userService, backgroundCtx)
	server.Start()
}

func mongoUri(config xrf.Config) (string, error) {
	mongoConfig := config.Database.Mongo
	baseUri := os.Getenv(mongoConfig.Uri)
	if baseUri == "" {
		baseUri = os.Getenv(mongoConfig.CloudUri)
	}
	if baseUri == "" {
		return "", &xrfErr.Internal{
			Source:  "cmd/cli/main#mongoUri",
			Message: "missing mongo uri environment variable $(uri/cloudUri)",
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
