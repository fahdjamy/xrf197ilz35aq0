package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"xrf197ilz35aq0"
)

type HealthRoutes struct {
	logger xrf197ilz35aq0.Logger
	router *mux.Router
}

func (hr *HealthRoutes) RegisterAndListen() {
	hr.router.HandleFunc("/health", hr.healthCheckHandler).Methods("GET")
}

func (hr *HealthRoutes) healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
	hr.logger.Info("api='/health' :: method='GET :: message='application health' :: status='ok'")
}

func NewHealthRoutes(logger xrf197ilz35aq0.Logger, router *mux.Router) *HealthRoutes {
	return &HealthRoutes{
		logger: logger,
		router: router,
	}
}
