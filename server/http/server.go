package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/core/service"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/server/http/handlers"
)

func RunServer(logger internal.Logger, config xrf197ilz35aq0.Config, services service.Services) {
	started := time.Now()

	//loggerMiddleware := middleware.NewLoggerHandler(server.logger)
	//authMiddleware := middleware.NewAuthenticationMiddleware(server.logger, server.services.AuthService)

	//server.router.Use(loggerMiddleware.Handler)

	// start the server
	appConfig := config.Application

	logger.Debug(fmt.Sprintf("timeouts :: readTO=%.2f :: writeTO=%.2f :: idleTO=%.2f :: graceShutdown=%.2f",
		appConfig.ReadTimeout.Seconds(),
		appConfig.WriteTimeout.Seconds(),
		appConfig.IdleTimeout.Seconds(),
		appConfig.GracefulTimeout.Seconds()))

	httpMux := http.NewServeMux()
	handlers.SetupHandlers(httpMux, logger, services)

	svr := http.Server{
		Handler:      httpMux,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
		IdleTimeout:  appConfig.IdleTimeout,
		Addr:         fmt.Sprintf(":%d", appConfig.Port),
	}

	// Run the server in a goroutine so that it doesn't block.
	go func() {
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error starting http server on port 8009: %s\n", err)
			logger.Error("serverStarted=false :: error starting http server")
		}
	}()

	timeTaken := time.Since(started).Milliseconds()
	logger.Info(fmt.Sprintf("serverStarted=true :: port=%d :: timeTaken='%d ms' message='application running...'", appConfig.Port, timeTaken))

	ch := make(chan os.Signal, 1)
	// Accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(ch, os.Interrupt)

	// Block until we receive shutdown signal.
	<-ch

	// Create a deadline context to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*appConfig.ReadTimeout)
	defer func() {
		cancel()
	}()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := svr.Shutdown(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("serverShutdown=failure :: %s", err))
		os.Exit(1)
	}

	// Improvement: Run svr.Shutdown in a goroutine and block on <-ctx.Done()
	// if your application should wait for other services to finalize based on context cancellation.

	logger.Info(fmt.Sprintf("serverShutdown=success"))
	os.Exit(0)
}
