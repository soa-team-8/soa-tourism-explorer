package handler

import (
	"encoding/json"
	"encounters/dto"
	"encounters/service"
	"encounters/utils"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

type HiddenLocationEncounterHandler struct {
	*utils.HttpUtils
	HiddenLocationEncounterService *service.HiddenLocationEncounterService
}

func NewHiddenLocationEncounterHandler(hiddenLocationEncounterService *service.HiddenLocationEncounterService) *HiddenLocationEncounterHandler {
	return &HiddenLocationEncounterHandler{
		HiddenLocationEncounterService: hiddenLocationEncounterService,
	}
}

func (s *HiddenLocationEncounterHandler) Create(resp http.ResponseWriter, req *http.Request) {
	newHiddenLocationEncounterDto := &dto.EncounterDto{}
	err := json.NewDecoder(req.Body).Decode(newHiddenLocationEncounterDto)
	if err != nil {
		s.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	err = req.ParseMultipartForm(10 << 20)
	if err != nil {
		s.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	images := req.MultipartForm.File["pictures"]

	imageService := service.NewImageService()
	uploadedImageNames, err := imageService.UploadImages(images)
	newHiddenLocationEncounterDto.Image = uploadedImageNames
	if err != nil {
		s.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	savedHiddenLocationEncounterDto, err := s.HiddenLocationEncounterService.Create(*newHiddenLocationEncounterDto)
	if err != nil {
		s.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	s.WriteJSONResponse(resp, http.StatusOK, savedHiddenLocationEncounterDto)
}

func (s *HiddenLocationEncounterHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	hiddenLocationEncounters, err := s.HiddenLocationEncounterService.GetAll()
	if err != nil {
		s.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	s.WriteJSONResponse(resp, http.StatusOK, hiddenLocationEncounters)
}

func (s *HiddenLocationEncounterHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := s.GetIDFromRequest(req)
	if err != nil {
		s.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	foundHiddenLocationEncounter, err := s.HiddenLocationEncounterService.GetByID(id)
	if err != nil {
		s.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	s.WriteJSONResponse(resp, http.StatusOK, foundHiddenLocationEncounter)
}

func (s *HiddenLocationEncounterHandler) UpdateByID(resp http.ResponseWriter, req *http.Request) {
	id, err := s.GetIDFromRequest(req)
	if err != nil {
		s.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedHiddenLocationEncounterDto := &dto.EncounterDto{}
	err = json.NewDecoder(req.Body).Decode(updatedHiddenLocationEncounterDto)
	if err != nil {
		s.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	updatedHiddenLocationEncounterDto.ID = id

	updatedHiddenLocationEncounter, err := s.HiddenLocationEncounterService.Update(*updatedHiddenLocationEncounterDto)
	if err != nil {
		s.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	s.WriteJSONResponse(resp, http.StatusOK, updatedHiddenLocationEncounter)
}

func (s *HiddenLocationEncounterHandler) DeleteByID(resp http.ResponseWriter, req *http.Request) {
	id, err := s.GetIDFromRequest(req)
	if err != nil {
		s.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	err = s.HiddenLocationEncounterService.DeleteByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.HandleError(resp, errors.New("social encounter not found"), http.StatusNotFound)
		} else {
			s.HandleError(resp, err, http.StatusInternalServerError)
		}
		return
	}

	s.WriteResponse(resp, http.StatusOK, "Social encounter deleted successfully")
}
