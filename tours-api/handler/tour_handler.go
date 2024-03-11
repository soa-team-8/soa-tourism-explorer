package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"tours/model"
	"tours/service"
	"tours/utils"
)

type TourHandler struct {
	TourService *service.TourService
	HttpUtils   *utils.HttpUtils
}

func (e *TourHandler) Create(resp http.ResponseWriter, req *http.Request) {
	tour, err := e.HttpUtils.Decode(req.Body, &model.Tour{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	id, err := e.TourService.Create(*tour.(*model.Tour))
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	tour.(*model.Tour).ID = id

	resp.WriteHeader(http.StatusCreated)
	resp.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(resp).Encode(tour)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}
}

func (e *TourHandler) Delete(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.TourService.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.HttpUtils.HandleError(resp, err, http.StatusNotFound)
		} else {
			e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		}
		return
	}

	e.HttpUtils.WriteResponse(resp, http.StatusOK, "Tour deleted successfully")
}

func (e *TourHandler) Update(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedTour, err := e.HttpUtils.Decode(req.Body, &model.Tour{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedTour.(*model.Tour).ID = id

	if err := e.TourService.Update(*updatedTour.(*model.Tour)); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteResponse(resp, http.StatusOK, "Tour updated successfully")
}

func (e *TourHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	equipment, err := e.TourService.GetByID(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	if equipment == nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("tour with ID %d not found", id), http.StatusNotFound)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, equipment)
}

func (e *TourHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	tour, err := e.TourService.GetAll()
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, tour)
}
