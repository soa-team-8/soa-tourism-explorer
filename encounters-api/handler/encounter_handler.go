package handler

import (
	"encoding/json"
	"encounters/dto"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"

	"errors"

	"gorm.io/gorm"

	"encounters/service"
	"encounters/utils"
)

type EncounterHandler struct {
	*utils.HttpUtils
	EncounterService *service.EncounterService
}

func NewEncounterHandler(encounterService *service.EncounterService) *EncounterHandler {
	return &EncounterHandler{
		EncounterService: encounterService,
	}
}

func (e *EncounterHandler) Create(resp http.ResponseWriter, req *http.Request) {
	newEncounterDto, err := e.Decode(req.Body, &dto.EncounterDto{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	savedEncounterDto, err := e.EncounterService.Create(*newEncounterDto.(*dto.EncounterDto))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, savedEncounterDto)
}

func (e *EncounterHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	encounters, err := e.EncounterService.GetAll()
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, encounters)
}

func (e *EncounterHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.GetIDFromRequest(req)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	foundEncounter, err := e.EncounterService.GetByID(id)
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, foundEncounter)
}

func (e *EncounterHandler) UpdateByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.GetIDFromRequest(req)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEncounterDto, err := e.Decode(req.Body, &dto.EncounterDto{})
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEncounterDto.(*dto.EncounterDto).ID = id

	updatedEncounter, err := e.EncounterService.Update(*updatedEncounterDto.(*dto.EncounterDto))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, updatedEncounter)
}

func (e *EncounterHandler) DeleteByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.GetIDFromRequest(req)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.EncounterService.DeleteByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.HandleError(resp, errors.New("encounter not found"), http.StatusNotFound)
		} else {
			e.HandleError(resp, err, http.StatusInternalServerError)
		}
		return
	}

	e.WriteResponse(resp, http.StatusOK, "Encounter deleted successfully")
}

func (e *EncounterHandler) CreateTouristEncounter(resp http.ResponseWriter, req *http.Request) {
	// Dobijanje vrednosti iz putanje pomoću gorilla/mux
	vars := mux.Vars(req)
	levelStr := vars["level"]
	userIDStr := vars["userId"]

	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	level, err := strconv.Atoi(levelStr)
	if err != nil {
		e.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	newEncounterDto := &dto.EncounterDto{}
	// Ostatak koda ostaje nepromenjen
	err = json.NewDecoder(strings.NewReader(req.FormValue("encounter"))).Decode(&newEncounterDto)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	images := req.MultipartForm.File["pictures"]

	imageService := service.NewImageService()
	uploadedImageNames, err := imageService.UploadImages(images)
	newEncounterDto.Image = uploadedImageNames
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	savedEncounterDto, err := e.EncounterService.CreateTouristEncounter(*newEncounterDto, level, uint64(userID))
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, savedEncounterDto)
}

func (e *EncounterHandler) CreateAuthorEncounter(resp http.ResponseWriter, req *http.Request) {
	// Dobijanje vrednosti iz putanje pomoću gorilla/mux

	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	newEncounterDto := &dto.EncounterDto{}
	// Ostatak koda ostaje nepromenjen
	err = json.NewDecoder(strings.NewReader(req.FormValue("encounter"))).Decode(&newEncounterDto)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	images := req.MultipartForm.File["pictures"]

	imageService := service.NewImageService()
	uploadedImageNames, err := imageService.UploadImages(images)
	newEncounterDto.Image = uploadedImageNames
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	savedEncounterDto, err := e.EncounterService.CreateAuthorEncounter(*newEncounterDto)
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, savedEncounterDto)
}
