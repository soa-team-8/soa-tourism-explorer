package handler

import (
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

	if err := e.TourService.Create(*tour.(*model.Tour)); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteResponse(resp, http.StatusCreated, "Tour created successfully")
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
