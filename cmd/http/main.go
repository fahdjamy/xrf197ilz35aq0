package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/cmd"
	"xrf197ilz35aq0/dependency"
	"xrf197ilz35aq0/server/http"
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

	router := mux.NewRouter().StrictSlash(true)

	apiServer := http.NewHttpServer(logger, router, config)
	apiServer.Start()
}
