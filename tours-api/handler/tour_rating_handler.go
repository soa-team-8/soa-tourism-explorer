package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"tours/model"
	"tours/service"
	"tours/utils"
)

type TourRatingHandler struct {
	TourRatingService *service.TourRatingService
	HttpUtils         *utils.HttpUtils
}

func (e *TourRatingHandler) Create(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	images := req.MultipartForm.File["images"]
	imageService := service.NewImageService()
	uploadedImageNames, err := imageService.UploadImages(images)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	createdRating := model.TourRating{}

	err = json.NewDecoder(strings.NewReader(req.FormValue("tourRating"))).Decode(&createdRating)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	createdRating.ImageNames = uploadedImageNames
	createdRating.ID = 0

	id, err := e.TourRatingService.Create(createdRating)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	createdRating.ID = id

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, createdRating)
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

	e.HttpUtils.WriteResponse(resp, http.StatusOK, "TourRating deleted successfully")
}

func (e *TourRatingHandler) Update(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	err = req.ParseMultipartForm(10 << 20)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	var updatedTourRating model.TourRating
	err = json.NewDecoder(strings.NewReader(req.FormValue("tourRating"))).Decode(&updatedTourRating)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedTourRating.ID = id

	if err = e.TourRatingService.Update(updatedTourRating); err != nil {
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
