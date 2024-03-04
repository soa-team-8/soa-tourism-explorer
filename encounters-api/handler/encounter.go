package handler

import (
	"fmt"
	"net/http"
)

type Encounter struct{}

func (e *Encounter) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create an encounter")
}

func (e *Encounter) GetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all encounters")
}

func (e *Encounter) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get an encounters by ID")
}

func (e *Encounter) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an encounters by ID")
}

func (e *Encounter) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an encounters by ID")
}
