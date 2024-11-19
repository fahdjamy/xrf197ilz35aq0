package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/cmd"
	"xrf197ilz35aq0/cmd/http/routes"
	"xrf197ilz35aq0/dependency"
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

	logPrefix := fmt.Sprintf("requestId='%s'", cmd.GenerateRequestId())
	logger := dependency.CustomZapLogger(environment.LogMode, config.Log.Level, logFileOutPut, logPrefix, initialFields)

	logger.Info(fmt.Sprintf("appVersion='%s' :: os='%s' :: message='application starting...'", health.Version(), health.Runtime.OS))

	router := mux.NewRouter().StrictSlash(true)

	api := routes.NewApi(logger, router, config)
	api.Start()
}
