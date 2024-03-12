package handler

import (
	"encounters/dto"
	"encounters/service"
	"encounters/utils"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type TouristEncounterHandler struct {
	*utils.HttpUtils
	EncounterService *service.EncounterService
}

func (e *TouristEncounterHandler) Create(resp http.ResponseWriter, req *http.Request) {
	newEncounterDto, err := e.Decode(req.Body, &dto.EncounterDto{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	savedEncounterDto, err := e.EncounterService.Create(*newEncounterDto.(*dto.EncounterDto))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, savedEncounterDto)
}

func (e *TouristEncounterHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	encounters, err := e.EncounterService.GetAll()
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, encounters)
}

func (e *TouristEncounterHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
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

func (e *TouristEncounterHandler) UpdateByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.GetIDFromRequest(req)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEncounterDto, err := e.Decode(req.Body, &dto.EncounterDto{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEncounterDto.(*dto.EncounterDto).ID = id

	updatedEncounter, err := e.EncounterService.Update(*updatedEncounterDto.(*dto.EncounterDto))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, updatedEncounter)
}

func (e *TouristEncounterHandler) DeleteByID(resp http.ResponseWriter, req *http.Request) {
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

func (e *TouristEncounterHandler) CreateTouristEncounter(resp http.ResponseWriter, req *http.Request) {
	// Dobijanje vrednosti iz putanje pomoću gorilla/mux
	vars := mux.Vars(req)
	checkpointIDStr := vars["checkpointId"]
	isSecretPrerequisiteStr := vars["isSecretPrerequisite"]
	levelStr := vars["level"]
	userIDStr := vars["userId"]

	// Pretvaranje stringova u odgovarajuće tipove
	checkpointID, err := strconv.Atoi(checkpointIDStr)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	isSecretPrerequisite, err := strconv.ParseBool(isSecretPrerequisiteStr)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	level, err := strconv.Atoi(levelStr)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	// Ostatak koda ostaje nepromenjen
	newEncounterDto, err := e.Decode(req.Body, &dto.EncounterDto{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	savedEncounterDto, err := e.EncounterService.CreateTouristEncounter(*newEncounterDto.(*dto.EncounterDto), checkpointID, isSecretPrerequisite, level, uint64(userID))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, savedEncounterDto)
}
