package handler

import (
	"fmt"
	"net/http"
)

type Tour struct{}

func (e *Tour) Create(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Create an tour")
}

func (e *Tour) GetAll(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get all tours")
}

func (e *Tour) GetByID(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get an tours by ID")
}

func (e *Tour) UpdateByID(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Update an tours by ID")
}

func (e *Tour) DeleteByID(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Delete an tours by ID")
}
