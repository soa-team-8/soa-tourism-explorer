package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/hello", basicHandler).Methods("GET")

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to listen to server", err)
	}
}

func basicHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello world"))
}
