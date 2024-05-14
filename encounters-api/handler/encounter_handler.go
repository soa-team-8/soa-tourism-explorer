package handler

import (
	"encounters/dto"
	"encounters/model"
	"errors"
	"mime/multipart"
	"net/http"

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
	id, err := e.GetUInt64FromRequest(req, "id")
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
	id, err := e.GetUInt64FromRequest(req, "id")
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
	id, err := e.GetUInt64FromRequest(req, "id")
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

func (e *EncounterHandler) CreateByTourist(resp http.ResponseWriter, req *http.Request) {
	level, err := e.HttpUtils.GetIntFromRequest(req, "level")
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	userID, err := e.HttpUtils.GetUInt64FromRequest(req, "userId")
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	newEncounterDto := &dto.EncounterDto{}
	err = e.HttpUtils.GetObjectFromForm(req, "encounter", newEncounterDto)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if newEncounterDto.Type == "Locaton" {
		images, err := e.HttpUtils.GetFilesFromForm(req, "pictures")
		if err != nil {
			e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
			return
		}
		err = e.uploadImages(resp, newEncounterDto, images)
		if err != nil {
			return
		}
	}

	if newEncounterDto.ActiveTouristsIDs == nil {
		newEncounterDto.ActiveTouristsIDs = &model.BigIntSlice{}
	}

	savedEncounterDto, err := e.EncounterService.CreateByTourist(*newEncounterDto, level, userID)
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, savedEncounterDto)
}

func (e *EncounterHandler) CreateByAuthor(resp http.ResponseWriter, req *http.Request) {
	newEncounterDto := &dto.EncounterDto{}
	err := e.HttpUtils.GetObjectFromForm(req, "encounter", newEncounterDto)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	if newEncounterDto.Type == "Locaton" {
		images, err := e.HttpUtils.GetFilesFromForm(req, "pictures")
		if err != nil {
			e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
			return
		}
		err = e.uploadImages(resp, newEncounterDto, images)
		if err != nil {
			return
		}
	}
	if newEncounterDto.ActiveTouristsIDs == nil {
		newEncounterDto.ActiveTouristsIDs = &model.BigIntSlice{}
	}
	savedEncounterDto, err := e.EncounterService.CreateByAuthor(*newEncounterDto)
	if err != nil {
		e.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.WriteJSONResponse(resp, http.StatusOK, savedEncounterDto)
}

func (e *EncounterHandler) uploadImages(resp http.ResponseWriter, newEncounterDto *dto.EncounterDto, images []*multipart.FileHeader) error {
	imageService := service.NewImageService()
	uploadedImageNames, err := imageService.UploadImages(images)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return err
	}
	newEncounterDto.Image = uploadedImageNames
	return nil
}
