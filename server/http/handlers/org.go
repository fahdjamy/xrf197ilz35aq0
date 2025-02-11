package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/internal/exchange"
	"xrf197ilz35aq0/server/http/decoder"
	"xrf197ilz35aq0/server/http/middleware"
	"xrf197ilz35aq0/server/http/response"
)

type OrgHandler struct {
	logger     xrf.Logger
	router     *mux.Router
	orgService service.OrgService
	authMiddle middleware.AuthenticationMiddleware
}

func NewOrgHandler(logger xrf.Logger, orgService service.OrgService, router *mux.Router, authMiddle middleware.AuthenticationMiddleware) *OrgHandler {
	return &OrgHandler{
		logger:     logger,
		router:     router,
		orgService: orgService,
		authMiddle: authMiddle,
	}
}

func (handler *OrgHandler) createOrg(w http.ResponseWriter, r *http.Request) {
	var orgReq exchange.OrgRequest
	err := decoder.DecodeJSONBody(r, &orgReq)
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}

	// create a new org
	resp, err := handler.orgService.CreateOrg(orgReq, context.Background())
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}
	dataResp := response.DataResponse{
		Code: 200,
		Data: struct {
			OrgId string `json:"orgId"`
		}{
			OrgId: resp,
		},
	}
	response.WriteResponse(dataResp, w, handler.logger)
}

func (handler *OrgHandler) updateOrg(w http.ResponseWriter, r *http.Request) {
	orgId, isValid := getAndValidateId(r, "orgId")
	if !isValid {
		externalError := &xrfErr.External{
			Message: "invalid org id",
		}
		response.WriteErrorResponse(externalError, w, handler.logger)
		return
	}
	var request exchange.UpdateOrgRequest
	err := decoder.DecodeJSONBody(r, &request)
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}

	// make call to update org
	resp, err := handler.orgService.UpdateOrg(orgId, request, context.Background())
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}

	dataResp := response.DataResponse{
		Code: 200,
		Data: struct {
			Updated bool                 `json:"updated"`
			Org     exchange.OrgResponse `json:"org"`
		}{
			Org:     *resp,
			Updated: resp != nil,
		},
	}
	response.WriteResponse(dataResp, w, handler.logger)
}

func (handler *OrgHandler) getOrg(w http.ResponseWriter, r *http.Request) {
	orgId, isValid := getAndValidateId(r, "orgId")
	if !isValid {
		externalError := &xrfErr.External{
			Code:    404,
			Message: "invalid org id",
		}
		response.WriteErrorResponse(externalError, w, handler.logger)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	foundOrg, err := handler.orgService.GetOrgById(orgId, ctx)
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}
	handler.logger.Debug(fmt.Sprintf("event=findOrg :: orgId=%s", orgId))

	resp := response.DataResponse{Data: foundOrg, Code: http.StatusOK}
	response.WriteResponse(resp, w, handler.logger)
}

func (handler *OrgHandler) findOrgMembers(w http.ResponseWriter, r *http.Request) {
	orgId, isValid := getAndValidateId(r, "orgId")
	if !isValid {
		externalError := &xrfErr.External{
			Message: "invalid org id",
		}
		response.WriteErrorResponse(externalError, w, handler.logger)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	foundOrgs, err := handler.orgService.FindOrgMembers(orgId, ctx)
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}
	handler.logger.Debug(fmt.Sprintf("event=findOrgMembers :: orgId=%s", orgId))
	resp := response.DataResponse{Data: foundOrgs, Code: http.StatusOK}
	response.WriteResponse(resp, w, handler.logger)
}

func (handler *OrgHandler) RegisterAndListen() {
	orgSubRoutes := handler.router.PathPrefix("/api/v1/org").Subrouter()

	orgSubRoutes.HandleFunc("", handler.createOrg).Methods(POST)
	orgSubRoutes.HandleFunc("/{orgId}", handler.getOrg).Methods(GET)
	orgSubRoutes.HandleFunc("/{orgId}", handler.updateOrg).Methods(PUT)
	orgSubRoutes.HandleFunc("/{orgId}/members", handler.findOrgMembers).Methods(GET)

	orgSubRoutes.Use(handler.authMiddle.Handle)
}
