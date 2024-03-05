package handler

import (
	"encoding/json"
	"encounters/model"
	"encounters/repository/encounter"
	"fmt"
	"net/http"
	"strconv"

	"errors"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Encounter struct {
	Repo *encounter.PostgresRepo
}

func (e *Encounter) Create(resp http.ResponseWriter, req *http.Request) {
	// Decode the JSON request body into an Encounter struct
	var newEncounter model.Encounter
	err := json.NewDecoder(req.Body).Decode(&newEncounter)
	if err != nil {
		http.Error(resp, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the Insert method of the repository to insert the new Encounter into the database
	err = e.Repo.Save(req.Context(), newEncounter)
	if err != nil {
		http.Error(resp, fmt.Sprintf("Failed to create encounter: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("Encounter created successfully"))
}

func (e *Encounter) GetAll(resp http.ResponseWriter, req *http.Request) {
	// Call the FindAll method of the repository to get all encounters
	encounters, err := e.Repo.FindAll(req.Context())
	if err != nil {
		http.Error(resp, fmt.Sprintf("Failed to get encounters: %v", err), http.StatusInternalServerError)
		return
	}

	// Encode the encounters into JSON
	encountersJSON, err := json.Marshal(encounters)
	if err != nil {
		http.Error(resp, "Failed to encode encounters", http.StatusInternalServerError)
		return
	}

	// Write the JSON response to the response writer
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(encountersJSON)
}

func (e *Encounter) GetByID(resp http.ResponseWriter, req *http.Request) {
	// Extract the encounter ID from the URL path parameters
	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(resp, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	foundEncounter, err := e.Repo.FindByID(req.Context(), id)
	if err != nil {
		http.Error(resp, fmt.Sprintf("Failed to get encounter: %v", err), http.StatusInternalServerError)
		return
	}

	encounterJSON, err := json.Marshal(foundEncounter)
	if err != nil {
		http.Error(resp, "Failed to encode encounter", http.StatusInternalServerError)
		return
	}

	// Write the JSON response to the response writer
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(encounterJSON)
}

func (e *Encounter) UpdateByID(resp http.ResponseWriter, req *http.Request) {
	// Extract the encounter ID from the URL path parameters
	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(resp, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	// Decode the JSON request body into an Encounter struct
	var updatedEncounter model.Encounter
	err = json.NewDecoder(req.Body).Decode(&updatedEncounter)
	if err != nil {
		http.Error(resp, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Set the ID of the updated encounter
	updatedEncounter.ID = id

	// Call the Update method of the repository to update the encounter
	err = e.Repo.Update(req.Context(), updatedEncounter)
	if err != nil {
		http.Error(resp, fmt.Sprintf("Failed to update encounter: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("Encounter updated successfully"))
}

func (e *Encounter) DeleteByID(resp http.ResponseWriter, req *http.Request) {
	// Extract the encounter ID from the URL path parameters
	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(resp, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	// Call the DeleteByID method of the repository to delete the encounter by its ID
	err = e.Repo.DeleteByID(req.Context(), id)

	if err != nil {
		// Check if the error is gorm.ErrRecordNotFound which indicates that the encounter was not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(resp, "Encounter not found", http.StatusNotFound)
		} else {
			http.Error(resp, fmt.Sprintf("Failed to delete encounter: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Respond with success message
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("Encounter deleted successfully"))
}

/*
func (e *Encounter) Create(resp http.ResponseWriter, req *http.Request) {
	// Define the JSON string representing the encounter
	jsonData := `{
		"author_id": 123,
		"id": 456,
		"name": "Exploration",
		"description": "An adventure in the wilderness",
		"XP": 100,
		"status": 2,
		"type": 1,
		"longitude": 45.6789,
		"latitude": 23.4567
	}`

	// Decode the JSON string into an Encounter struct
	var newEncounter model.Encounter
	err := json.Unmarshal([]byte(jsonData), &newEncounter)
	if err != nil {
		http.Error(resp, "Failed to decode encounter JSON", http.StatusInternalServerError)
		return
	}

	// Call the Insert method of the repository to insert the new Encounter into the database
	err = e.Repo.Save(req.Context(), newEncounter)
	if err != nil {
		http.Error(resp, fmt.Sprintf("Failed to create encounter: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("Encounter created successfully"))
}
*/
