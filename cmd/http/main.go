package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
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

	defer func(mongoClient *mongo2.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			logger.Error(fmt.Sprintf("appStarted=false :: err%s", err.Error()))
		}
	}(mongoClient, backgroundCtx)
	if err != nil {
		internalError := xrfErr.Internal{Err: err, Message: "failed to connect to mongo"}
		logger.Error(fmt.Sprintf("appStarted=failure :: %s", internalError.Error()))
		return
	}
	// connect to mongoDB
	mongoDB := mongoClient.Database(databaseName)
	logger.Debug(fmt.Sprintf("message='successfully connected to MongoDB' :: dbName=%s", databaseName))

	// create repositories
	permissionRepo, err := repository.NewPermissionRepo(mongoDB, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("appStarted=false :: err%s", err.Error()))
		return
	}
	orgRepo, err := repository.NewOrganizationRepository(mongoDB, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("appStarted=false :: err%s", err.Error()))
		return
	}

	userRepo := repository.NewUserRepository(mongoDB, logger)
	settingRepo := repository.NewSettingsRepository(mongoDB, logger)

	allRepos := &repository.Repositories{
		PermissionRepo: permissionRepo,
		UserRepo:       userRepo,
		OrgRepo:        orgRepo,
		SettingsRepo:   settingRepo,
	}

	// create services
	permService := service.NewPermissionService(logger, permissionRepo)
	orgService := service.NewOrganizationService(config.Security, logger, allRepos)
	settingsService := service.NewSettingService(logger, settingRepo, backgroundCtx, config.Security)
	userService := service.NewUserService(logger, settingsService, userRepo, backgroundCtx, config.Security)

	services := http.Services{
		OrgService:        orgService,
		UserService:       userService,
		PermissionService: permService,
	}

	// create the router and start the server
	router := mux.NewRouter().StrictSlash(true)
	server := http.NewHttpServer(logger, router, config, services, backgroundCtx)
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
