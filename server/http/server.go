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
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/server/http/handlers"
	"xrf197ilz35aq0/storage"
)

type ApiServer struct {
	started  bool
	logger   internal.Logger
	ctx      context.Context
	services service.Services
	cache    storage.Cache
	config   xrf197ilz35aq0.Config
}

var apiInternalErr = &xrfErr.Internal{
	Source: "cmd/http/run#start",
}

func (server *ApiServer) Start() {
	if server.started {
		return
	}
	started := time.Now()

	//loggerMiddleware := middleware.NewLoggerHandler(server.logger)
	//authMiddleware := middleware.NewAuthenticationMiddleware(server.logger, server.services.AuthService)

	// handlers
	//handlers.NewAuthHandler(server.logger, server.services, server.router).RegisterAndListen()
	//handlers.NewUserHandler(server.logger, server.services.UserService, server.router).RegisterAndListen()
	//handlers.NewPermHandler(server.logger, server.router, server.services.PermissionService).RegisterAndListen()
	//handlers.NewOrgHandler(server.logger, server.services.OrgService, server.router, authMiddleware).RegisterAndListen()

	//server.router.Use(loggerMiddleware.Handler)

	// start the server
	appConfig := server.config.Application

	server.logger.Debug(fmt.Sprintf("timeouts :: readTO=%.2f :: writeTO=%.2f :: idleTO=%.2f :: graceShutdown=%.2f",
		appConfig.ReadTimeout.Seconds(),
		appConfig.WriteTimeout.Seconds(),
		appConfig.IdleTimeout.Seconds(),
		appConfig.GracefulTimeout.Seconds()))

	svr := http.Server{
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
			server.logger.Error(fmt.Sprintf("serverStarted=false :: %s", apiInternalErr))
		}
	}()

	timeTaken := time.Since(started).Milliseconds()
	server.logger.Info(fmt.Sprintf("serverStarted=true :: port=%d :: timeTaken='%d ms' message='application running...'", appConfig.Port, timeTaken))
	server.started = true

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
		server.Stop()
	}()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := svr.Shutdown(ctx)
	if err != nil {
		server.logger.Error(fmt.Sprintf("serverShutdown=failure :: %s", err))
		os.Exit(1)
	}

	// Improvement: Run svr.Shutdown in a goroutine and block on <-ctx.Done()
	// if your application should wait for other services to finalize based on context cancellation.

	server.logger.Info(fmt.Sprintf("serverShutdown=success"))
	os.Exit(0)
}

func (server *ApiServer) Stop() {
	if !server.started {
		return
	}
	server.started = false
}

func SetUpHttpServer(logger internal.Logger, config xrf197ilz35aq0.Config, services service.Services, ctx context.Context) http.Handler {
	httpMux := http.NewServeMux()
	routes := make([]RoutesHandler, 0)

	routes = append(routes, handlers.NewHealthRoutes(logger))
	SetUpRoutes(httpMux, routes)

	var handler http.Handler = httpMux

	return handler
}
