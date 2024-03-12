package application

import (
	"encounters/handler"
	"encounters/repo"
	"encounters/service"

	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (a *App) loadRoutes() {
	router := mux.NewRouter()
	router.Use(loggerMiddleware)

	router.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
	}).Methods("GET")

	encountersRouter := router.PathPrefix("/encounters").Subrouter()
	a.loadEncounterRoutes(encountersRouter)

	a.router = router
}

func (a *App) loadEncounterRoutes(router *mux.Router) {
	encounterService := &service.EncounterService{
		EncounterRepo: &repo.EncounterRepository{
			DB: a.db,
		},
	}

	encounterHandler := &handler.EncounterHandler{
		EncounterService: encounterService,
	}

	touristEncounterHandler := &handler.TouristEncounterHandler{
		EncounterService: encounterService,
	}

	router.HandleFunc("", encounterHandler.Create).Methods("POST")
	router.HandleFunc("", encounterHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", encounterHandler.GetByID).Methods("GET")
	router.HandleFunc("/{id}", encounterHandler.UpdateByID).Methods("PUT")
	router.HandleFunc("/{id}", encounterHandler.DeleteByID).Methods("DELETE")

	router.HandleFunc("/tourist", touristEncounterHandler.Create).Methods("POST")
	router.HandleFunc("/tourist", touristEncounterHandler.GetAll).Methods("GET")
	router.HandleFunc("/tourist/{id}", touristEncounterHandler.GetByID).Methods("GET")
	router.HandleFunc("/tourist/{id}", touristEncounterHandler.UpdateByID).Methods("PUT")
	router.HandleFunc("/tourist/{id}", touristEncounterHandler.DeleteByID).Methods("DELETE")
	router.HandleFunc("/tourist/{checkpointId}/{isSecretPrerequisite}/{level}/{userId}/", touristEncounterHandler.CreateTouristEncounter).Methods("POST")
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
