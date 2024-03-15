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

type TourRatingHandler struct {
	TourRatingService *service.TourRatingService
	HttpUtils         *utils.HttpUtils
}

func (e *TourRatingHandler) Create(resp http.ResponseWriter, req *http.Request) {
	tourRating, err := e.HttpUtils.Decode(req.Body, &model.TourRating{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.TourRatingService.Create(*tourRating.(*model.TourRating)); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, tourRating)
}

func (e *TourRatingHandler) Delete(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.TourRatingService.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.HttpUtils.HandleError(resp, err, http.StatusNotFound)
		} else {
			e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		}
		return
	}

	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteResponse(resp, http.StatusOK, "Tour deleted successfully")
}

func (e *TourRatingHandler) Update(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedTourRating, err := e.HttpUtils.Decode(req.Body, &model.TourRating{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedTourRating.(*model.TourRating).ID = id

	if err := e.TourRatingService.Update(*updatedTourRating.(*model.TourRating)); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, updatedTourRating)
}

func (e *TourRatingHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	tourRating, err := e.TourRatingService.GetByID(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	if tourRating == nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("tourRating with ID %d not found", id), http.StatusNotFound)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, tourRating)
}

func (e *TourRatingHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	tourRating, err := e.TourRatingService.GetAll()
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, tourRating)
}
