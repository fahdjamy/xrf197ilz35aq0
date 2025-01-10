package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
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

func (handler *OrgHandler) getOrg(w http.ResponseWriter, r *http.Request) {
	orgId, isValid := getAndValidateId(r, "orgId")
	if !isValid {
		externalError := &xrfErr.External{
			Message: "invalid org id",
		}
		writeErrorResponse(externalError, w, handler.logger)
		return
	}
	ctx, close := context.WithTimeout(context.Background(), time.Second*2)
	defer close()
	foundOrg, err := handler.orgService.GetOrgById(orgId, ctx)
	if err != nil {
		writeErrorResponse(err, w, handler.logger)
		return
	}
	handler.logger.Debug(fmt.Sprintf("event=findOrg :: orgId=%s", orgId))

	resp := dataResponse{Data: foundOrg, Code: http.StatusOK}
	writeResponse(resp, w, handler.logger)
}

func (handler *OrgHandler) RegisterAndListen() {
	slashAPISlashOrg := fmt.Sprintf("%s/%s/%s", constants.SlashAPI, constants.V1, "org") // "/api/v1/org"
	findByOrgIdUrl := fmt.Sprintf("%s/{%s}", slashAPISlashOrg, constants.OrgId)          // "/api/v1/org/{orgId}"

	handler.router.HandleFunc(findByOrgIdUrl, handler.getOrg).Methods(GET)
	handler.router.HandleFunc(slashAPISlashOrg, handler.createOrg).Methods(POST)
}
