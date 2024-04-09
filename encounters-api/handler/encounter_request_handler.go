package handler

import (
	"encounters/dto"
	"encounters/service"
	"encounters/utils"
	"net/http"
)

type EncounterRequestHandler struct {
	*utils.HttpUtils
	EncounterRequestService *service.EncounterRequestService
}

func NewEncounterRequestHandler(executionRequestService *service.EncounterRequestService) *EncounterRequestHandler {
	return &EncounterRequestHandler{
		EncounterRequestService: executionRequestService,
	}
}

func (e *EncounterRequestHandler) CreateRequest(resp http.ResponseWriter, req *http.Request) {
	encounterReqDto, err := e.Decode(req.Body, &dto.EncounterRequestDto{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	createdReq, err := e.EncounterRequestService.CreateEncounterRequest(*encounterReqDto.(*dto.EncounterRequestDto))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusCreated, createdReq)
}

func (e *EncounterRequestHandler) GetRequestByID(resp http.ResponseWriter, req *http.Request) {
	idReq, err := e.GetIDFromRequest(req, "id")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	encounterRequest, err := e.EncounterRequestService.GetEncounterRequestByID(int(idReq))
	if err != nil {
		e.HandleError(resp, err, http.StatusNotFound)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, encounterRequest)
}

func (e *EncounterRequestHandler) UpdateRequest(resp http.ResponseWriter, req *http.Request) {
	updatedEncounterRequestDto, err := e.Decode(req.Body, &dto.EncounterRequestDto{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEncounter, err := e.EncounterRequestService.Update(*updatedEncounterRequestDto.(*dto.EncounterRequestDto))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, updatedEncounter)
}

func (e *EncounterRequestHandler) DeleteRequest(resp http.ResponseWriter, req *http.Request) {
	idReq, err := e.GetIDFromRequest(req, "id")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	err = e.EncounterRequestService.DeleteByID(int(idReq))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteResponse(resp, http.StatusOK, "Encounter deleted successfully")
}

func (e *EncounterRequestHandler) AcceptRequest(resp http.ResponseWriter, req *http.Request) {
	idReq, err := e.GetIDFromRequest(req, "id")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	acceptedReq, err := e.EncounterRequestService.Accept(int(idReq))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, acceptedReq)
}

func (e *EncounterRequestHandler) RejectRequest(resp http.ResponseWriter, req *http.Request) {
	idReq, err := e.GetIDFromRequest(req, "id")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	rejectedReq, err := e.EncounterRequestService.Reject(int(idReq))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, rejectedReq)
}

func (e *EncounterRequestHandler) GetAllRequests(resp http.ResponseWriter, req *http.Request) {
	encounterRequests, err := e.EncounterRequestService.GetAll()
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, encounterRequests)
}
