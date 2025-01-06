package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
)

type RoleHandler struct {
	logger      xrf.Logger
	router      *mux.Router
	roleService service.RoleService
}

func (handler *RoleHandler) createRole(w http.ResponseWriter, r *http.Request) {
	var roleReq *exchange.RoleRequest
	err := decodeJSONBody(r, &roleReq)
	if err != nil {
		writeErrorResponse(err, w, handler.logger)
		return
	}

	// create a new role
	resp, err := handler.roleService.CreateRole(roleReq, context.Background())
	if err != nil {
		writeErrorResponse(err, w, handler.logger)
		return
	}
	dataResp := dataResponse{
		Code: 200,
		Data: struct {
			RoleId string `json:"roleId"`
		}{
			RoleId: resp,
		},
	}
	writeResponse(dataResp, w, handler.logger)
}

func (handler *RoleHandler) RegisterAndListen() {
	handler.router.HandleFunc("/role", handler.createRole).Methods("POST")
}

func NewRoleHandler(logger xrf.Logger, router *mux.Router, service service.RoleService) *RoleHandler {
	return &RoleHandler{
		logger:      logger,
		router:      router,
		roleService: service,
	}
}
