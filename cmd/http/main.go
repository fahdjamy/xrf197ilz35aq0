package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strconv"
	"xrf197ilz35aq0/internal/constants"
	"xrf197ilz35aq0/internal/random"
	"xrf197ilz35aq0/storage"

	"github.com/redis/go-redis/v9"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/core/repository"
	"xrf197ilz35aq0/core/service"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/dependency"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/server/http"
	"xrf197ilz35aq0/storage/mongo"
)

const (
	AuthSecretEnvKey = "XRF_Q0_AUTH_SECRET_KEY"
)

func main() {
	// get the globally set environment variables
	environment := internal.GetEnvironment()
	serverId := random.PositiveInt64()

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
	logPrefix := fmt.Sprintf("requestId='%s'", fmt.Sprintf("%d||starting", serverId))
	logger := dependency.CustomZapLogger(environment.LogMode, config.Log.Level, logFileOutPut, logPrefix, initialFields)
	logger.Info(fmt.Sprintf("appVersion='%s' :: os='%s' :: message='application starting...' :: serverId=%d",
		health.Version(),
		health.Runtime.OS,
		serverId),
	)

	authSecret, exists := os.LookupEnv(AuthSecretEnvKey)
	if !exists {
		logger.Error(fmt.Sprintf("appStarted=false :: message='Missing _Q0_AUTH_SECRET_KEY'"))
		return
	}

	// connect to redis server
	redisClient, err := connectRedis(config.Redis, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("appStarted=false :: message='Could not start redis' :: %v", err.Error()))
		return
	}

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

	userRepo, err := repository.NewUserRepository(mongoDB, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("appStarted=false :: err%s", err.Error()))
		return
	}

	settingRepo := repository.NewSettingsRepository(mongoDB, logger)

	allRepos := &repository.Repositories{
		PermissionRepo: permissionRepo,
		UserRepo:       userRepo,
		OrgRepo:        orgRepo,
		SettingsRepo:   settingRepo,
	}

	// set cache service
	redisCache := storage.NewRedisStorage(redisClient)

	// create services
	permService := service.NewPermissionService(logger, permissionRepo)
	orgService := service.NewOrganizationService(config.Security, logger, allRepos)
	settingsService := service.NewSettingService(logger, settingRepo, backgroundCtx, config.Security)
	userService := service.NewUserService(logger, settingsService, userRepo, backgroundCtx, config.Security)
	authService := service.NewAuthService(strconv.FormatInt(serverId, 10), logger, authSecret, redisCache, allRepos)

	services := service.Services{
		OrgService:        orgService,
		UserService:       userService,
		AuthService:       authService,
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

func connectRedis(config xrf.RedisConfig, logger internal.Logger) (*redis.Client, error) {
	redisAddress := config.Address

	if redisAddress == "" {
		logger.Error(fmt.Sprintf("event=connectRedis :: message='looking for redis address in environment..."))
		redisAddressInEnv, ok := os.LookupEnv(constants.RedisAddress)
		redisAddress = redisAddressInEnv
		if !ok {
			return &redis.Client{}, &xrfErr.Internal{
				Source:  "cmd/cli/main#redis",
				Message: "missing redis address environment variable $(redis_address)",
			}
		}
	}
	client := redis.NewClient(&redis.Options{
		Addr:       config.Address,
		Password:   config.Password,
		DB:         config.Database,
		Protocol:   config.Protocol,
		MaxRetries: config.MaxRetries,
	})

	// Ping the Redis server to check the connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("event=connectRedis :: action=redisServerStarted :: address=%s", config.Address))

	return client, nil
}
