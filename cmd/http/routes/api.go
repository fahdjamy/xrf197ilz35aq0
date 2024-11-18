package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"xrf197ilz35aq0"
	xrfErr "xrf197ilz35aq0/internal/error"
)

type Route struct {
	started bool
	router  *mux.Router
	logger  xrf197ilz35aq0.Logger
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
	err := http.ListenAndServe(":8009", r.router)
	if err != nil {
		fmt.Printf("Error starting http server on port 8009: %s\n", err)
		apiInternalErr.Time = time.Now()
		apiInternalErr.Message = "error starting http server"
		r.logger.Error(fmt.Sprintf("serverStarted=false :: %s", apiInternalErr))
		//return
	}

	timeTaken := time.Since(started)
	r.logger.Info(fmt.Sprintf("serverStarted=true :: time=%s :: timeTake=%s :: message=http server started", started, timeTaken))
	r.started = true
}

func (r *Route) Stop() {
	if !r.started {
		return
	}
	r.started = false
}

func NewApi(logger xrf197ilz35aq0.Logger, router *mux.Router) *Route {
	return &Route{
		logger: logger,
		router: router,
	}
}
