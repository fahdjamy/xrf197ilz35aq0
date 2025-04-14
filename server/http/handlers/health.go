package handlers

import (
	"fmt"
	"net/http"
	xrf "xrf197ilz35aq0/internal"
)

type HealthRoutes struct {
	logger xrf.Logger
}

func (hr *HealthRoutes) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		hr.logger.Error(fmt.Sprintf("event=healthCheckFailure :: message='Setting header failed' :: err=%s", err.Error()))
		return
	}
}

func (hr *HealthRoutes) RegisterAndListen() {
}

func (hr *HealthRoutes) RegisterRoutes(serveMux *http.ServeMux) {
	serveMux.Handle("GET /health", http.Handler(http.HandlerFunc(hr.healthCheck)))
}

func NewHealthRoutes(logger xrf.Logger) *HealthRoutes {
	return &HealthRoutes{
		logger: logger,
	}
}
