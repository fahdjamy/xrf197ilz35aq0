package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
)

type PermissionHandler struct {
	logger      xrf.Logger
	router      *mux.Router
	permService service.PermissionService
}

func (handler *PermissionHandler) createPermission(w http.ResponseWriter, r *http.Request) {
	var permissionReq *exchange.PermissionRequest
	err := decodeJSONBody(r, &permissionReq)
	if err != nil {
		writeErrorResponse(err, w, handler.logger)
		return
	}

	// create a new permission
	resp, err := handler.permService.CreatePermission(permissionReq, context.Background())
	if err != nil {
		writeErrorResponse(err, w, handler.logger)
		return
	}
	dataResp := dataResponse{
		Code: 200,
		Data: struct {
			Id string `json:"id"`
		}{
			Id: resp,
		},
	}
	writeResponse(dataResp, w, handler.logger)
}

func (handler *PermissionHandler) RegisterAndListen() {
	handler.router.HandleFunc("/permission", handler.createPermission).Methods("POST")
}

func NewPermHandler(logger xrf.Logger, router *mux.Router, service service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		logger:      logger,
		router:      router,
		permService: service,
	}
}
