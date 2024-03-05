package application

import (
	"fmt"
	"net/http"
	"time"
	"tours/handler"

	"github.com/gorilla/mux"
)

func loadRoutes() *mux.Router {
	router := mux.NewRouter()
	router.Use(loggerMiddleware)

	router.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
	}).Methods("GET")

	toursRouter := router.PathPrefix("/tours").Subrouter()
	loadTourRoutes(toursRouter)

	return router
}

func loadTourRoutes(router *mux.Router) {
	toursRouter := &handler.Tour{}

	router.HandleFunc("/", toursRouter.Create).Methods("POST")
	router.HandleFunc("/", toursRouter.GetAll).Methods("GET")
	router.HandleFunc("/{id}", toursRouter.GetByID).Methods("GET")
	router.HandleFunc("/{id}", toursRouter.UpdateByID).Methods("PUT")
	router.HandleFunc("/{id}", toursRouter.DeleteByID).Methods("DELETE")
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()

		next.ServeHTTP(resp, req)

		fmt.Printf(
			"[%s] %s %s %v\n",
			req.Method,
			req.RequestURI,
			req.RemoteAddr,
			time.Since(start),
		)
	})
}
