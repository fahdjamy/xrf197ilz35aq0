package routes

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
	"xrf197ilz35aq0"
	xrfErr "xrf197ilz35aq0/internal/error"
)

type Route struct {
	started bool
	router  *mux.Router
	logger  xrf197ilz35aq0.Logger
	config  xrf197ilz35aq0.Config
}

var apiInternalErr = &xrfErr.Internal{
	Source: "cmd/http/run#start",
}

func (r *Route) Start() {
	if r.started {
		return
	}
	started := time.Now()

	// handlers
	NewHealthRoutes(r.logger, r.router).RegisterAndListen()

	// start the server
	appConfig := r.config.Application
	r.logger.Debug(fmt.Sprintf("timeouts :: readTO=%d :: writeTO=%d :: idleTO=%d :: graceShutdown=%d",
		time.Second*appConfig.ReadTimeout,
		time.Second*appConfig.WriteTimeout,
		appConfig.IdleTimeoutSecs,
		appConfig.GracefulTimeoutSecs))

	svr := http.Server{
		ReadTimeout:  time.Second * appConfig.ReadTimeout,
		WriteTimeout: time.Second * appConfig.WriteTimeout,
		IdleTimeout:  time.Second * appConfig.IdleTimeoutSecs,
		Handler:      r.router,
		Addr:         fmt.Sprintf(":%d", appConfig.Port),
	}

	// Run the server in a goroutine so that it doesn't block.
	go func() {
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error starting http server on port 8009: %s\n", err)
			apiInternalErr.Time = time.Now()
			apiInternalErr.Message = "error starting http server"
			r.logger.Error(fmt.Sprintf("serverStarted=false :: %s", apiInternalErr))
		}
	}()

	timeTaken := time.Since(started).Milliseconds()
	r.logger.Info(fmt.Sprintf("serverStarted=true :: port=%d :: timeTaken='%d ms'", appConfig.Port, timeTaken))
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
		r.started = false
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

func (r *Route) Stop() {
	if !r.started {
		return
	}
	r.started = false
}

func NewApi(logger xrf197ilz35aq0.Logger, router *mux.Router, config xrf197ilz35aq0.Config) *Route {
	return &Route{
		logger: logger,
		router: router,
		config: config,
	}
}
