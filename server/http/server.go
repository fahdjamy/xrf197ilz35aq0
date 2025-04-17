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

func RunServer(logger internal.Logger, config xrf197ilz35aq0.ApplicationConfig, svr *http.Server) {
	started := time.Now()

	// Run the server in a goroutine so that it doesn't block.
	go func() {
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error starting http server on port 8009: %s\n", err)
			logger.Error("serverStarted=false :: error starting http server")
		}
	}()

	timeTaken := time.Since(started).Milliseconds()
	logger.Info(fmt.Sprintf("serverStarted=true :: port=%d :: timeTaken='%d ms' message='application running...'", config.Port, timeTaken))

	ch := make(chan os.Signal, 1)
	// Accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(ch, os.Interrupt)

	// Block until we receive shutdown signal.
	<-ch

	// Create a deadline context to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*config.ReadTimeout)
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

func CreateServer(logger internal.Logger, services service.Services, config xrf197ilz35aq0.ApplicationConfig) *http.Server {
	idleTimeout := config.IdleTimeout
	readTimeout := config.ReadTimeout
	writeTimeout := config.WriteTimeout
	gracefulTimeout := config.GracefulTimeout

	logger.Debug(fmt.Sprintf("timeouts :: readTO=%.2f :: writeTO=%.2f :: idleTO=%.2f :: graceShutdown=%.2f",
		readTimeout.Seconds(),
		writeTimeout.Seconds(),
		idleTimeout.Seconds(),
		gracefulTimeout.Seconds()))

	serverMux := http.NewServeMux()

	// 2. register handlers
	wrappedServer := handlers.SetupHandlers(serverMux, logger, services)

	return &http.Server{
		Handler:      wrappedServer,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
		Addr:         fmt.Sprintf(":%d", config.Port),
	}
}
