package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/core/service/user"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/server/http/handlers"
	"xrf197ilz35aq0/server/http/middleware"
)

type ApiServer struct {
	started     bool
	router      *mux.Router
	logger      xrf.Logger
	config      xrf.Config
	userManager user.Manager
	ctx         context.Context
}

var apiInternalErr = &xrfErr.Internal{
	Source: "cmd/http/run#start",
}

func (r *ApiServer) Start() {
	if r.started {
		return
	}
	started := time.Now()

	loggerMiddleware := middleware.NewLoggerHandler(r.logger)

	// handlers
	r.router.Use(loggerMiddleware.Handler)
	handlers.NewHealthRoutes(r.logger, r.router).RegisterAndListen()
	handlers.NewUser(r.logger, r.userManager, r.router).RegisterAndListen()

	// start the server
	appConfig := r.config.Application

	r.logger.Debug(fmt.Sprintf("timeouts :: readTO=%.2f :: writeTO=%.2f :: idleTO=%.2f :: graceShutdown=%.2f",
		appConfig.ReadTimeout.Seconds(),
		appConfig.WriteTimeout.Seconds(),
		appConfig.IdleTimeout.Seconds(),
		appConfig.GracefulTimeout.Seconds()))

	svr := http.Server{
		Handler:      r.router,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
		IdleTimeout:  appConfig.IdleTimeout,
		Addr:         fmt.Sprintf(":%d", appConfig.Port),
	}

	// Run the server in a goroutine so that it doesn't block.
	go func() {
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error starting http server on port 8009: %s\n", err)
			apiInternalErr.Message = "error starting http server"
			r.logger.Error(fmt.Sprintf("serverStarted=false :: %s", apiInternalErr))
		}
	}()

	timeTaken := time.Since(started).Milliseconds()
	r.logger.Info(fmt.Sprintf("serverStarted=true :: port=%d :: timeTaken='%d ms' message='application running...'", appConfig.Port, timeTaken))
	r.started = true

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
		r.Stop()
	}()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := svr.Shutdown(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("serverShutdown=failure :: %s", err))
		os.Exit(1)
	}

	// Improvement: Run svr.Shutdown in a goroutine and block on <-ctx.Done()
	// if your application should wait for other services to finalize based on context cancellation.

	r.logger.Info(fmt.Sprintf("serverShutdown=success"))
	os.Exit(0)
}

func (r *ApiServer) Stop() {
	if !r.started {
		return
	}
	r.started = false
}

func NewHttpServer(logger xrf.Logger, router *mux.Router, config xrf.Config, userManager user.Manager, ctx context.Context) *ApiServer {
	return &ApiServer{
		ctx:         ctx,
		logger:      logger,
		router:      router,
		config:      config,
		userManager: userManager,
	}
}
