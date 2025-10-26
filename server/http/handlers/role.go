package handlers

import (
	"net/http"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/exchange"
	"xrf197ilz35aq0/server/http/decoder"
	"xrf197ilz35aq0/server/http/response"
)

type PermissionHandler struct {
	logger      xrf.Logger
	permService service.PermissionService
}

func (handler *PermissionHandler) createPermission(w http.ResponseWriter, r *http.Request) {
	var permissionReq *exchange.PermissionRequest
	err := decoder.DecodeJSONBody(r, &permissionReq)
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}

	// create a new permission
	resp, err := handler.permService.CreatePermission(r.Context(), permissionReq)
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}
	dataResp := response.DataResponse{
		Code: 200,
		Data: struct {
			Id string `json:"id"`
		}{
			Id: resp,
		},
	}
	response.WriteResponse(dataResp, w, handler.logger)
}

func (handler *PermissionHandler) RegisterRoutes(serveMux *http.ServeMux) {
	serveMux.HandleFunc("POST /api/v1/permission", handler.createPermission)
}

func NewPermHandler(logger xrf.Logger, service service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		logger:      logger,
		permService: service,
	}
}
