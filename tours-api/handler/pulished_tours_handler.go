package handler

import (
	"net/http"
	"tours/service"
	"tours/utils"
)

type PublishedToursHandler struct {
	TourService *service.TourService
	HttpUtils   *utils.HttpUtils
}

func (e *PublishedToursHandler) GetPublishedTours(resp http.ResponseWriter, req *http.Request) {
	tourPreviews, err := e.TourService.GetPublishedTours()
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}
	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, tourPreviews)
}

func (e *PublishedToursHandler) GetPublishedTour(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	tourPreview, err := e.TourService.GetPublishedTour(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}
	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, tourPreview)
}

func (e *PublishedToursHandler) GetAverageRating(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	rating, err := e.TourService.GetAverageRating(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}
	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, rating)
}
