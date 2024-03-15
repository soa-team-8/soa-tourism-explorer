package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"tours/model"
	"tours/service"
	"tours/utils"
)

type CheckpointHandler struct {
	CheckpointService *service.CheckpointService
	HttpUtils         *utils.HttpUtils
}

func (e *CheckpointHandler) Create(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	images := req.MultipartForm.File["pictures"]

	imageService := service.NewImageService()
	uploadedImageNames, err := imageService.UploadImages(images)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	createdCheckpoint := model.Checkpoint{}

	err = json.NewDecoder(strings.NewReader(req.FormValue("checkpoint"))).Decode(&createdCheckpoint)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	createdCheckpoint.Pictures = uploadedImageNames

	id, err := e.CheckpointService.Create(createdCheckpoint)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	createdCheckpoint.ID = id

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, createdCheckpoint)
}

func (e *CheckpointHandler) Update(resp http.ResponseWriter, req *http.Request) {
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

	images := req.MultipartForm.File["pictures"]
	var uploadedImageNames []string
	if len(images) > 0 {
		imageService := service.NewImageService()
		uploadedImageNames, err = imageService.UploadImages(images)
		if err != nil {
			e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
			return
		}
	}

	var updatedCheckpoint model.Checkpoint
	err = json.NewDecoder(strings.NewReader(req.FormValue("checkpoint"))).Decode(&updatedCheckpoint)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if len(uploadedImageNames) > 0 {
		updatedCheckpoint.Pictures = uploadedImageNames
	}

	updatedCheckpoint.ID = id

	if err := e.CheckpointService.Update(updatedCheckpoint); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, updatedCheckpoint)
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

func (e *CheckpointHandler) CreateOrUpdateCheckpointSecret(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	images := req.MultipartForm.File["pictures"]
	uploadedImageNames, err := service.NewImageService().UploadImages(images)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	var checkpointSecret model.CheckpointSecret
	err = json.NewDecoder(strings.NewReader(req.FormValue("checkpointSecret"))).Decode(&checkpointSecret)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	checkpointSecret.Pictures = uploadedImageNames

	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.CheckpointService.CreateOrUpdateCheckpointSecret(id, checkpointSecret); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	updatedCheckpoint, err := e.CheckpointService.GetByID(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, updatedCheckpoint)
}

func (e *CheckpointHandler) SetCheckpointEncounter(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// Izvlačenje parametara iz putanje
	checkpointID, err := strconv.ParseInt(vars["checkpointId"], 10, 64)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	encounterID, err := strconv.ParseInt(vars["encId"], 10, 64)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	isSecretPrerequisite, err := strconv.ParseBool(vars["isSecretPrerequisite"])
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	// Poziv metode iz servisa
	err = e.CheckpointService.SetCheckpointEncounter(uint64(checkpointID), uint64(encounterID), isSecretPrerequisite)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, "Checkpoint encounter successfully set")
}
