package handler

import (
	"encounters/model"
	"encounters/service"
	"encounters/utils"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

type EncounterExecutionHandler struct {
	*utils.HttpUtils
	ExecutionService *service.EncounterExecutionService
}

func NewEncounterExecutionHandler(executionService *service.EncounterExecutionService) *EncounterExecutionHandler {
	return &EncounterExecutionHandler{
		ExecutionService: executionService,
	}
}

func (e *EncounterExecutionHandler) Create(resp http.ResponseWriter, req *http.Request) {
	touristId, err := e.GetIDFromRequest(req, "touristId")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	newExecution, err := e.Decode(req.Body, &model.EncounterExecution{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	savedEncounter, err := e.ExecutionService.Create(*newExecution.(*model.EncounterExecution), touristId)
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, savedEncounter)
}

func (e *EncounterExecutionHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	encounters, err := e.ExecutionService.GetAll()
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, encounters)
}

func (e *EncounterExecutionHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.GetIDFromRequest(req, "id")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	foundEncounter, err := e.ExecutionService.GetByID(id)
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, foundEncounter)
}

func (e *EncounterExecutionHandler) DeleteByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.GetIDFromRequest(req, "id")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	touristId, err := e.GetIDFromRequest(req, "touristId")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.ExecutionService.DeleteByID(id, touristId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.HandleError(resp, errors.New("execution not found"), http.StatusNotFound)
		} else {
			e.HandleError(resp, err, http.StatusInternalServerError)
		}
		return
	}

	e.WriteResponse(resp, http.StatusOK, "Execution deleted successfully")
}

func (e *EncounterExecutionHandler) UpdateByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.GetIDFromRequest(req, "id")
	touristId, err := e.GetIDFromRequest(req, "touristId")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedExecution, err := e.Decode(req.Body, &model.EncounterExecution{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedExecution.(*model.EncounterExecution).ID = id

	updatedEncounter, err := e.ExecutionService.Update(*updatedExecution.(*model.EncounterExecution), touristId)
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, updatedEncounter)
}
