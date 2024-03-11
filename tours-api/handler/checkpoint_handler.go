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

type CheckpointHandler struct {
	CheckpointService *service.CheckpointService
	HttpUtils         *utils.HttpUtils
}

func (e *CheckpointHandler) Create(resp http.ResponseWriter, req *http.Request) {
	checkpoint, err := e.HttpUtils.Decode(req.Body, &model.Checkpoint{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.CheckpointService.Create(*checkpoint.(*model.Checkpoint)); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteResponse(resp, http.StatusCreated, "Checkpoint created successfully")
}

func (e *CheckpointHandler) Delete(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.CheckpointService.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.HttpUtils.HandleError(resp, err, http.StatusNotFound)
		} else {
			e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		}
		return
	}

	e.HttpUtils.WriteResponse(resp, http.StatusOK, "Checkpoint deleted successfully")
}

func (e *CheckpointHandler) Update(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedCheckpoint, err := e.HttpUtils.Decode(req.Body, &model.Checkpoint{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedCheckpoint.(*model.Checkpoint).ID = id

	if err := e.CheckpointService.Update(*updatedCheckpoint.(*model.Checkpoint)); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteResponse(resp, http.StatusOK, "Checkpoint updated successfully")
}

func (e *CheckpointHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	checkpoint, err := e.CheckpointService.GetByID(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	if checkpoint == nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("checkpoint with ID %d not found", id), http.StatusNotFound)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, checkpoint)
}

func (e *CheckpointHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	checkpoint, err := e.CheckpointService.GetAll()
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, checkpoint)
}

func (e *CheckpointHandler) GetAllByTourID(resp http.ResponseWriter, req *http.Request) {
	tourId, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	checkpoints, err := e.CheckpointService.GetAllByTourID(tourId)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	if checkpoints == nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("checkpoints for Tour with ID %d not found", tourId), http.StatusNotFound)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, checkpoints)
}
