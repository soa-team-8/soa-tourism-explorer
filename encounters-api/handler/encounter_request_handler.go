package handler

import (
	"encounters/service"
	"encounters/utils"
	"net/http"
)

type EncounterRequestHandler struct {
	*utils.HttpUtils
	EncounterRequestService *service.EncounterRequestService
}

func (e *EncounterRequestHandler) AcceptRequest(resp http.ResponseWriter, req *http.Request) {
	idReq, err := e.GetIDFromRequest(req)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	acceptedReq, err := e.EncounterRequestService.AcceptEncounterRequest(int(idReq))

	e.WriteJSONResponse(resp, http.StatusOK, acceptedReq)
}

func (e *EncounterRequestHandler) RejectRequest(resp http.ResponseWriter, req *http.Request) {
	idReq, err := e.GetIDFromRequest(req)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	acceptedReq, err := e.EncounterRequestService.RejectEncounterRequest(int(idReq))

	e.WriteJSONResponse(resp, http.StatusOK, acceptedReq)
}

func (e *EncounterRequestHandler) GetAllRequests(resp http.ResponseWriter, req *http.Request) {
	encounterRequests, err := e.EncounterRequestService.GetAll()
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, encounterRequests)
}
