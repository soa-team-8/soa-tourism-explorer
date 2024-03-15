package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tours/model"
	"tours/service"
	"tours/utils"
)

type TourExecutionHandler struct {
	TourExecutionService *service.TourExecutionService
	HttpUtils            *utils.HttpUtils
}

func (e *TourExecutionHandler) CheckPosition(resp http.ResponseWriter, req *http.Request) {
	touristPosition, err := e.HttpUtils.Decode(req.Body, &model.TouristPosition{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if tourExecution, err := e.TourExecutionService.CheckPosition(*touristPosition.(*model.TouristPosition), id); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	} else {
		e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, tourExecution)
	}
}

func (e *TourExecutionHandler) Abandon(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userIdStr := vars["userID"]
	executionIdStr := vars["executionID"]

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	executionID, err := strconv.Atoi(executionIdStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if tourExecution, err := e.TourExecutionService.Abandon(uint64(userID), uint64(executionID)); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	} else {
		e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, tourExecution)
	}
}

func (e *TourExecutionHandler) Create(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userIdStr := vars["userID"]
	tourIdStr := vars["tourID"]

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	tourID, err := strconv.Atoi(tourIdStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if tourExecution, err := e.TourExecutionService.Create(uint64(userID), uint64(tourID)); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	} else {
		e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, tourExecution)
	}
}

func (e *TourExecutionHandler) GetByIDs(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userIdStr := vars["userID"]
	tourIdStr := vars["tourID"]

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	tourID, err := strconv.Atoi(tourIdStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	tourExecution, err := e.TourExecutionService.GetByIDs(uint64(userID), uint64(tourID))
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	if tourExecution == nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("tourExecution with IDs %d, %d not found", userID, tourID), http.StatusNotFound)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, tourExecution)
}
