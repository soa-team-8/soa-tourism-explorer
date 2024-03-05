package handler

import (
	"fmt"
	"net/http"
)

type Encounter struct{}

func (e *Encounter) Create(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Create an encounter")
}

func (e *Encounter) GetAll(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get all encounters")
}

func (e *Encounter) GetByID(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get an encounters by ID")
}

func (e *Encounter) UpdateByID(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Update an encounters by ID")
}

func (e *Encounter) DeleteByID(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Delete an encounters by ID")
}
