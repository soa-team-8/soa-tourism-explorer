package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tours/model"
	"tours/service"
	"tours/utils"

	"gorm.io/gorm"
)

type EquipmentHandler struct {
	EquipmentService *service.EquipmentService
	HttpUtils        *utils.HttpUtils
}

func (e *EquipmentHandler) Create(resp http.ResponseWriter, req *http.Request) {
	equipment, err := e.HttpUtils.Decode(req.Body, &model.Equipment{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.EquipmentService.Create(*equipment.(*model.Equipment)); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, equipment)
}

func (e *EquipmentHandler) Delete(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.EquipmentService.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.HttpUtils.HandleError(resp, err, http.StatusNotFound)
		} else {
			e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		}
		return
	}

	responseData := map[string]string{"message": "Equipment deleted successfully"}
	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(jsonResponse)
}

func (e *EquipmentHandler) Update(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEquipment, err := e.HttpUtils.Decode(req.Body, &model.Equipment{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEquipment.(*model.Equipment).ID = id

	if err := e.EquipmentService.Update(*updatedEquipment.(*model.Equipment)); err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, updatedEquipment)
}

func (e *EquipmentHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	equipment, err := e.EquipmentService.GetByID(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	if equipment == nil {
		e.HttpUtils.HandleError(resp, fmt.Errorf("equipment with ID %d not found", id), http.StatusNotFound)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, equipment)
}

func (e *EquipmentHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	equipment, err := e.EquipmentService.GetAll()
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, equipment)
}

func (e *EquipmentHandler) GetAllPaged(resp http.ResponseWriter, req *http.Request) {
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	pageSize, err := strconv.Atoi(req.URL.Query().Get("pageSize"))
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	equipment, err := e.EquipmentService.GetAllPaged(page, pageSize)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, equipment)
}
