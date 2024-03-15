package handler

import (
	"encounters/model"
	"encounters/service"
	"encounters/utils"
	"errors"
	"fmt"
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
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

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

func (e *EncounterExecutionHandler) Activate(resp http.ResponseWriter, req *http.Request) {
	touristID, err := e.GetIDFromRequest(req, "touristId")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	executionID, err := e.GetIDFromRequest(req, "encounterId")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	// Parse parameters from form data
	touristLongitude, err := e.HttpUtils.GetDoubleFromForm(req, "touristLongitude")
	if err != nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("invalid tourist longitude"), http.StatusBadRequest)
		return
	}

	touristLatitude, err := e.HttpUtils.GetDoubleFromForm(req, "touristLatitude")
	if err != nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("invalid tourist latitude"), http.StatusBadRequest)
		return
	}

	updatedExecution, err := e.ExecutionService.Activate(executionID, touristID, touristLongitude, touristLatitude)
	if err != nil {
		http.Error(resp, fmt.Sprintf("Error activating encounter execution: %v", err), http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, updatedExecution)
}

func (e *EncounterExecutionHandler) Complete(resp http.ResponseWriter, req *http.Request) {
	touristID, err := e.HttpUtils.GetIDFromRequest(req, "touristId")
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	executionID, err := e.HttpUtils.GetIDFromRequest(req, "executionId")
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	// Parse longitude and latitude from form data
	touristLongitude, err := e.HttpUtils.GetDoubleFromForm(req, "touristLongitude")
	if err != nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("invalid tourist longitude"), http.StatusBadRequest)
		return
	}

	touristLatitude, err := e.HttpUtils.GetDoubleFromForm(req, "touristLatitude")
	if err != nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("invalid tourist latitude"), http.StatusBadRequest)
		return
	}

	updatedExecution, XP, err := e.ExecutionService.Complete(executionID, touristID, touristLongitude, touristLatitude)
	if err != nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("error completing encounter execution: %v", err), http.StatusInternalServerError)
		return
	}

	responseData := struct {
		Execution *model.EncounterExecution `json:"execution"`
		XP        int32                     `json:"xp"`
	}{
		Execution: updatedExecution,
		XP:        XP,
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, responseData)
}

func (e *EncounterExecutionHandler) GetByTour(resp http.ResponseWriter, req *http.Request) {
	touristId, err := e.GetIDFromRequest(req, "touristId")
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	encounters, err := e.ExecutionService.GetAllByTourist(touristId)
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, encounters)

}

func (e *EncounterExecutionHandler) GetActiveByTour(resp http.ResponseWriter, req *http.Request) {

}

func (e *EncounterExecutionHandler) GetAllByTourist(resp http.ResponseWriter, req *http.Request) {

}

func (e *EncounterExecutionHandler) GetAllCompletedByTourist(resp http.ResponseWriter, req *http.Request) {

}

func (e *EncounterExecutionHandler) CheckPosition(resp http.ResponseWriter, req *http.Request) {

}

func (e *EncounterExecutionHandler) CheckPositionLocationEncounter(resp http.ResponseWriter, req *http.Request) {

}
