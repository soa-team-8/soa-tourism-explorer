package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"tours/model"
	"tours/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type EquipmentHandler struct {
	EquipmentService *service.EquipmentService
}

func (e *EquipmentHandler) Create(resp http.ResponseWriter, req *http.Request) {
	newEquipment, err := decodeEquipment(req.Body)
	if err != nil {
		handleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.EquipmentService.Create(*newEquipment); err != nil {
		handleError(resp, fmt.Errorf("failed to create equipment: %w", err), http.StatusInternalServerError)
		return
	}

	writeResponse(resp, http.StatusCreated, "Equipment created successfully")
}

func (e *EquipmentHandler) Delete(resp http.ResponseWriter, req *http.Request) {
	id, err := getEquipmentIDFromRequest(req)
	if err != nil {
		handleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.EquipmentService.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			handleError(resp, errors.New("equipment not found"), http.StatusNotFound)
		} else {
			handleError(resp, fmt.Errorf("failed to delete equipment: %w", err), http.StatusInternalServerError)
		}
		return
	}

	writeResponse(resp, http.StatusOK, "Equipment deleted successfully")
}

func (e *EquipmentHandler) Update(resp http.ResponseWriter, req *http.Request) {
	id, err := getEquipmentIDFromRequest(req)
	if err != nil {
		handleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEquipment, err := decodeEquipment(req.Body)
	if err != nil {
		handleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEquipment.ID = id

	if err := e.EquipmentService.Update(*updatedEquipment); err != nil {
		handleError(resp, fmt.Errorf("failed to update equipment: %w", err), http.StatusInternalServerError)
		return
	}

	writeResponse(resp, http.StatusOK, "Equipment updated successfully")
}

func (e *EquipmentHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := getEquipmentIDFromRequest(req)
	if err != nil {
		handleError(resp, err, http.StatusBadRequest)
		return
	}

	equipment, err := e.EquipmentService.GetByID(id)
	if err != nil {
		handleError(resp, fmt.Errorf("failed to get equipment: %w", err), http.StatusInternalServerError)
		return
	}

	if equipment == nil {
		handleError(resp, fmt.Errorf("equipment with ID %d not found", id), http.StatusNotFound)
		return
	}

	writeJSONResponse(resp, http.StatusOK, equipment)
}

func (e *EquipmentHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	equipment, err := e.EquipmentService.GetAll()
	if err != nil {
		handleError(resp, fmt.Errorf("failed to get all equipment: %w", err), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(resp, http.StatusOK, equipment)
}

func decodeEquipment(body io.Reader) (*model.Equipment, error) {
	var equipment model.Equipment
	err := json.NewDecoder(body).Decode(&equipment)
	return &equipment, err
}

func getEquipmentIDFromRequest(req *http.Request) (uint64, error) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	return strconv.ParseUint(idStr, 10, 64)
}

func handleError(resp http.ResponseWriter, err error, statusCode int) {
	http.Error(resp, fmt.Sprintf("Error: %v", err), statusCode)
}

func writeJSONResponse(resp http.ResponseWriter, statusCode int, data interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	err := json.NewEncoder(resp).Encode(data)
	if err != nil {
		return
	}
}

func writeResponse(resp http.ResponseWriter, statusCode int, message string) {
	resp.WriteHeader(statusCode)
	resp.Write([]byte(message))
}
