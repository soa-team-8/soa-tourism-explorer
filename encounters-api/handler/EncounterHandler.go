package handler

import (
	"encoding/json"
	"encounters/model"
	"fmt"
	"net/http"
	"strconv"

	"errors"

	"io"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"encounters/service"
)

type EncounterHandler struct {
	EncounterService *service.EncounterService
}

func (e *EncounterHandler) Create(resp http.ResponseWriter, req *http.Request) {
	newEncounter, err := decodeEncounter(req.Body)
	if err != nil {
		handleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.EncounterService.Create(*newEncounter); err != nil {
		handleError(resp, err, http.StatusInternalServerError)
		return
	}

	writeResponse(resp, http.StatusCreated, "Encounter created successfully")
}

func (e *EncounterHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	encounters, err := e.EncounterService.GetAll()
	if err != nil {
		handleError(resp, err, http.StatusInternalServerError)
		return
	}

	writeJSONResponse(resp, http.StatusOK, encounters)
}

func (e *EncounterHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	id, err := getEncounterIDFromRequest(req)
	if err != nil {
		handleError(resp, err, http.StatusBadRequest)
		return
	}

	foundEncounter, err := e.EncounterService.GetByID(id)
	if err != nil {
		handleError(resp, err, http.StatusInternalServerError)
		return
	}

	writeJSONResponse(resp, http.StatusOK, foundEncounter)
}

func (e *EncounterHandler) UpdateByID(resp http.ResponseWriter, req *http.Request) {
	id, err := getEncounterIDFromRequest(req)
	if err != nil {
		handleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEncounter, err := decodeEncounter(req.Body)
	if err != nil {
		handleError(resp, err, http.StatusBadRequest)
		return
	}

	updatedEncounter.ID = id

	if err := e.EncounterService.Update(*updatedEncounter); err != nil {
		handleError(resp, err, http.StatusInternalServerError)
		return
	}

	writeResponse(resp, http.StatusOK, "Encounter updated successfully")
}

func (e *EncounterHandler) DeleteByID(resp http.ResponseWriter, req *http.Request) {
	id, err := getEncounterIDFromRequest(req)
	if err != nil {
		handleError(resp, err, http.StatusBadRequest)
		return
	}

	if err := e.EncounterService.DeleteByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			handleError(resp, errors.New("encounter not found"), http.StatusNotFound)
		} else {
			handleError(resp, err, http.StatusInternalServerError)
		}
		return
	}

	writeResponse(resp, http.StatusOK, "Encounter deleted successfully")
}

func decodeEncounter(body io.Reader) (*model.Encounter, error) {
	var newEncounter model.Encounter
	err := json.NewDecoder(body).Decode(&newEncounter)
	return &newEncounter, err
}

func getEncounterIDFromRequest(req *http.Request) (uint64, error) {
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
	json.NewEncoder(resp).Encode(data)
}

func writeResponse(resp http.ResponseWriter, statusCode int, message string) {
	resp.WriteHeader(statusCode)
	resp.Write([]byte(message))
}

/*
{
	"author_id": 123,
	"id": 456,
	"name": "Exploration",
	"description": "An adventure in the wilderness",
	"XP": 100,
	"status": 2,
	"type": 1,
	"longitude": 45.6789,
	"latitude": 23.4567
}
*/
