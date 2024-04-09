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

type SocialEncounterHandler struct {
	*utils.HttpUtils
	SocialEncounterService *service.SocialEncounterService
}

func NewSocialEncounterHandler(socialEncounterService *service.SocialEncounterService) *SocialEncounterHandler {
	return &SocialEncounterHandler{
		SocialEncounterService: socialEncounterService,
	}
}

func (s *SocialEncounterHandler) Create(resp http.ResponseWriter, req *http.Request) {
	newSocialEncounterDto := &dto.EncounterDto{}
	err := json.NewDecoder(req.Body).Decode(newSocialEncounterDto)
	if err != nil {
		s.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	savedSocialEncounterDto, err := s.SocialEncounterService.Create(*newSocialEncounterDto)
	if err != nil {
		s.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	s.WriteJSONResponse(resp, http.StatusOK, savedSocialEncounterDto)
}

func (s *SocialEncounterHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	socialEncounters, err := s.SocialEncounterService.GetAll()
	if err != nil {
		s.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	s.WriteJSONResponse(resp, http.StatusOK, socialEncounters)
}

func (s *SocialEncounterHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := s.GetUInt64FromRequest(req, "id")
	if err != nil {
		s.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	foundSocialEncounter, err := s.SocialEncounterService.GetByID(id)
	if err != nil {
		s.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	s.WriteJSONResponse(resp, http.StatusOK, foundSocialEncounter)
}

func (s *SocialEncounterHandler) UpdateByID(resp http.ResponseWriter, req *http.Request) {
	id, err := s.GetUInt64FromRequest(req, "id")
	if err != nil {
		s.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedSocialEncounterDto := &dto.EncounterDto{}
	err = json.NewDecoder(req.Body).Decode(updatedSocialEncounterDto)
	if err != nil {
		s.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	updatedSocialEncounterDto.ID = id

	updatedSocialEncounter, err := s.SocialEncounterService.Update(*updatedSocialEncounterDto)
	if err != nil {
		s.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	s.WriteJSONResponse(resp, http.StatusOK, updatedSocialEncounter)
}

func (s *SocialEncounterHandler) DeleteByID(resp http.ResponseWriter, req *http.Request) {
	id, err := s.GetUInt64FromRequest(req, "id")
	if err != nil {
		s.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	err = s.SocialEncounterService.DeleteByID(id)
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
