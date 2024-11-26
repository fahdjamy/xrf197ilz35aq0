package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	xrf "xrf197ilz35aq0/internal"
)

type HealthRoutes struct {
	logger xrf.Logger
	router *mux.Router
}

func (hr *HealthRoutes) RegisterAndListen() {
	hr.router.HandleFunc("/health", hr.healthCheck).Methods("GET")
}

func (hr *HealthRoutes) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		hr.logger.Error(fmt.Sprintf("event=healthCheckFailure :: message='Setting header failed' :: err=%s", err.Error()))
		return
	}
}

func NewHealthRoutes(logger xrf.Logger, router *mux.Router) *HealthRoutes {
	return &HealthRoutes{
		logger: logger,
		router: router,
	}
}
