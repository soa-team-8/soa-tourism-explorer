package handler

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, tour)
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

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, updatedTour)
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

func (e *TourHandler) AddEquipmentToTour(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tourID, err := strconv.ParseUint(vars["tourID"], 10, 64)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	equipmentID, err := strconv.ParseUint(vars["equipmentID"], 10, 64)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.TourService.AddEquipmentToTour(tourID, equipmentID); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	updatedTour, err := e.TourService.GetByID(tourID)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, updatedTour)
}

func (e *TourHandler) RemoveEquipmentFromTour(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tourID, err := strconv.ParseUint(vars["tourID"], 10, 64)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	equipmentID, err := strconv.ParseUint(vars["equipmentID"], 10, 64)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.TourService.RemoveEquipmentFromTour(tourID, equipmentID); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	updatedTour, err := e.TourService.GetByID(tourID)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, updatedTour)
}

func (e *TourHandler) GetToursByAuthor(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	authorIDStr := vars["authorID"]

	authorID, err := strconv.ParseUint(authorIDStr, 10, 64)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	tours, err := e.TourService.GetToursByAuthor(authorID)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, tours)
}

func (e *TourHandler) Publish(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tourID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	tour, err := e.TourService.GetByID(tourID)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	if tour == nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("tour with ID %d not found", tourID), http.StatusNotFound)
		return
	}

	if len(tour.Checkpoints) >= 2 {
		tour.Status = 1

		if err := e.TourService.Update(*tour); err != nil {
			e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
			return
		}

		e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, tour)
	} else {
		e.HttpUtils.HandleError(resp, errors.New("tour must have at least two checkpoints to be published"), http.StatusBadRequest)
	}
}

func (e *TourHandler) Archive(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tourID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	tour, err := e.TourService.GetByID(tourID)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	if tour == nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("tour with ID %d not found", tourID), http.StatusNotFound)
		return
	}

	if tour.Status == 1 {
		tour.Status = 2

		if err := e.TourService.Update(*tour); err != nil {
			e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
			return
		}

		e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, tour)
	} else {
		e.HttpUtils.HandleError(resp, errors.New("tour must be published to be archived"), http.StatusBadRequest)
	}
}
