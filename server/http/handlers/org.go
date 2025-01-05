package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
)

type OrgHandler struct {
	logger     xrf.Logger
	router     *mux.Router
	orgService service.OrgService
}

func NewOrgHandler(logger xrf.Logger, orgService service.OrgService, router *mux.Router) *OrgHandler {
	return &OrgHandler{
		logger:     logger,
		router:     router,
		orgService: orgService,
	}
}

func (handler *OrgHandler) createOrg(w http.ResponseWriter, r *http.Request) {
	var orgReq exchange.OrgRequest
	err := decodeJSONBody(r, &orgReq)
	if err != nil {
		writeErrorResponse(err, w, handler.logger)
		return
	}

	// create a new org
	resp, err := handler.orgService.CreateOrg(orgReq, context.Background())
	if err != nil {
		writeErrorResponse(err, w, handler.logger)
		return
	}
	dataResp := dataResponse{
		Code: 200,
		Data: struct {
			OrgId string `json:"orgId"`
		}{
			OrgId: resp,
		},
	}
	writeResponse(dataResp, w, handler.logger)
}

func (handler *OrgHandler) RegisterAndListen() {
	handler.router.HandleFunc("/org", handler.createOrg).Methods("POST")
}
