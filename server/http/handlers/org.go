package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/internal/exchange"
	"xrf197ilz35aq0/server/http/decoder"
	"xrf197ilz35aq0/server/http/response"
)

type OrgHandler struct {
	logger     xrf.Logger
	orgService service.OrgService
}

func NewOrgHandler(logger xrf.Logger, orgService service.OrgService) *OrgHandler {
	return &OrgHandler{
		logger:     logger,
		orgService: orgService,
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
	resp, err := handler.orgService.CreateOrg(r.Context(), orgReq)
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
	resp, err := handler.orgService.UpdateOrg(r.Context(), orgId, request)
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
		externalError := &xrfErr.External{Code: 404, Message: "invalid org id"}
		response.WriteErrorResponse(externalError, w, handler.logger)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()
	foundOrg, err := handler.orgService.GetOrgById(ctx, orgId)
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}
	handler.logger.Debug(fmt.Sprintf("event=findOrg :: orgId=%s", orgId))

	resp := response.DataResponse{Data: foundOrg, Code: http.StatusOK}
	response.WriteResponse(resp, w, handler.logger)
}

func (handler *OrgHandler) getOrgOrDefault(w http.ResponseWriter, r *http.Request) {
	orgId, isValid := getAndValidateId(r, "orgId")
	if !isValid {
		externalError := &xrfErr.External{Code: 404, Message: "invalid org id"}
		response.WriteErrorResponse(externalError, w, handler.logger)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()
	foundOrg, err := handler.orgService.GetOrgOrDefault(ctx, orgId)
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}
	handler.logger.Debug(fmt.Sprintf("event=getOrgOrDefault :: orgId=%s", orgId))
	resp := response.DataResponse{Data: foundOrg, Code: http.StatusOK}
	response.WriteResponse(resp, w, handler.logger)
}

func (handler *OrgHandler) findOrgMembers(w http.ResponseWriter, r *http.Request) {
	orgId, isValid := getAndValidateId(r, "orgId")
	if !isValid {
		externalError := &xrfErr.External{Message: "invalid org id"}
		response.WriteErrorResponse(externalError, w, handler.logger)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	foundOrgs, err := handler.orgService.FindOrgMembers(ctx, orgId)
	if err != nil {
		response.WriteErrorResponse(err, w, handler.logger)
		return
	}
	handler.logger.Debug(fmt.Sprintf("event=findOrgMembers :: orgId=%s", orgId))
	resp := response.DataResponse{Data: foundOrgs, Code: http.StatusOK}
	response.WriteResponse(resp, w, handler.logger)
}

func (handler *OrgHandler) RegisterRoutes(serveMux *http.ServeMux) {
	serveMux.HandleFunc("POST /api/v1/org", handler.createOrg)
	serveMux.HandleFunc("GET /api/v1/org/{orgId}", handler.getOrg)
	serveMux.HandleFunc("PUT /api/v1/org/{orgId}", handler.updateOrg)
	serveMux.HandleFunc("GET /api/v1/org/{orgId}/members", handler.findOrgMembers)
	serveMux.HandleFunc("GET /api/v1/org-or-default/{orgId}", handler.getOrgOrDefault)
}
