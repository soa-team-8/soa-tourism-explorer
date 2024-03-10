package handler

import (
	"encounters/dto"
	"net/http"

	"errors"

	"gorm.io/gorm"

	"encounters/service"
	"encounters/utils"
)

type EncounterHandler struct {
	*utils.HttpUtils // Embedding the HttpUtils struct
	EncounterService *service.EncounterService
}

func NewEncounterHandler(httpUtils *utils.HttpUtils, service *service.EncounterService) *EncounterHandler {
	return &EncounterHandler{
		HttpUtils:        httpUtils,
		EncounterService: service,
	}
}

func (e *EncounterHandler) Create(resp http.ResponseWriter, req *http.Request) {
	newEncounter, err := e.Decode(req.Body, &dto.EncounterDto{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.EncounterService.Create(*newEncounter.(*dto.EncounterDto)); err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteResponse(resp, http.StatusCreated, "Encounter created successfully")
}

func (e *EncounterHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	encounters, err := e.EncounterService.GetAll()
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, encounters)
}

func (e *EncounterHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.GetIDFromRequest(req)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	foundEncounter, err := e.EncounterService.GetByID(id)
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, foundEncounter)
}

func (e *EncounterHandler) UpdateByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.GetIDFromRequest(req)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEncounter, err := e.Decode(req.Body, &dto.EncounterDto{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEncounter.(*dto.EncounterDto).ID = id

	if err := e.EncounterService.Update(*updatedEncounter.(*dto.EncounterDto)); err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteResponse(resp, http.StatusOK, "Encounter updated successfully")
}

func (e *EncounterHandler) DeleteByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.GetIDFromRequest(req)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.EncounterService.DeleteByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.HandleError(resp, errors.New("encounter not found"), http.StatusNotFound)
		} else {
			e.HandleError(resp, err, http.StatusInternalServerError)
		}
		return
	}

	e.WriteResponse(resp, http.StatusOK, "Encounter deleted successfully")
}

/*
{
	"author_id": 123,
	"id": 456,
	"name": "Exploration",
	"description": "An adventure in the wilderness",
	"XP": 100,
	"status": 2,
	"type": 1,
	"longitude": 45.6789,
	"latitude": 23.4567
}
*/
