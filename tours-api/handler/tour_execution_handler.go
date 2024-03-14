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
	
}

func (e *TourExecutionHandler) Abandon(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uidStr := vars["uid"]
	eidStr := vars["eid"]

	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	eid, err := strconv.Atoi(eidStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	var tourExecution model.TourExecution
	if _, err := e.TourExecutionService.Abandon(uid, eid); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, tourExecution) //nema nista te
}

func (e *TourExecutionHandler) Create(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uidStr := vars["uid"]
	tidStr := vars["tid"]

	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	tid, err := strconv.Atoi(tidStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	var tourExecution model.TourExecution
	if _, err := e.TourExecutionService.Create(uid, tid); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, tourExecution) //nema nista te
}

func (e *TourExecutionHandler) GetByIDs(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uidStr := vars["uid"]
	tidStr := vars["tid"]

	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	tid, err := strconv.Atoi(tidStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	tourExecution, err := e.TourExecutionService.GetByIDs(uid, tid)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	if tourExecution == nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("tourExecution with IDs %d, %d not found", uid, tid), http.StatusNotFound)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, tourExecution)
}
