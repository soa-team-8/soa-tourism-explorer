package application

import (
	"encounters/handler"
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
	executionsRouter := router.PathPrefix("/encounters/execution").Subrouter()

	a.loadEncounterRoutes(encountersRouter)
	a.loadExecutionRoutes(executionsRouter)

	encounterRequestRouter := router.PathPrefix("/requests").Subrouter()
	a.loadEncounterRequestRoutes(encounterRequestRouter)

	socialEncounterRouter := router.PathPrefix("/encounters/social").Subrouter()
	a.loadSocialEncounterRoutes(socialEncounterRouter)

	a.router = router
}

func (a *App) loadEncounterRoutes(router *mux.Router) {
	encounterService := service.NewEncounterService(a.db)
	encounterHandler := handler.NewEncounterHandler(encounterService)

	router.HandleFunc("", encounterHandler.Create).Methods("POST")
	router.HandleFunc("/get-all", encounterHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", encounterHandler.GetByID).Methods("GET")
	router.HandleFunc("/{id}", encounterHandler.UpdateByID).Methods("PUT")
	router.HandleFunc("/{id}", encounterHandler.DeleteByID).Methods("DELETE")
	router.HandleFunc("/tourist/{checkpointId}/{isSecretPrerequisite}/{level}/{userId}", encounterHandler.CreateTouristEncounter).Methods("POST")
}

func (a *App) loadEncounterRequestRoutes(router *mux.Router) {
	encounterRequestService := service.NewEncounterRequestService(a.db)
	encounterRequestHandler := handler.NewEncounterRequestHandler(encounterRequestService)

	router.HandleFunc("/create", encounterRequestHandler.CreateRequest).Methods("POST")
	router.HandleFunc("/get/{id}", encounterRequestHandler.GetRequestByID).Methods("GET")
	router.HandleFunc("/update", encounterRequestHandler.UpdateRequest).Methods("PUT")
	router.HandleFunc("/delete/{id}", encounterRequestHandler.DeleteRequest).Methods("DELETE")
	router.HandleFunc("/acceptReq/{id}", encounterRequestHandler.AcceptRequest).Methods("PUT")
	router.HandleFunc("/rejectReq/{id}", encounterRequestHandler.RejectRequest).Methods("PUT")
	router.HandleFunc("/getAll", encounterRequestHandler.GetAllRequests).Methods("GET")
}

func (a *App) loadExecutionRoutes(router *mux.Router) {
	executionService := service.NewEncounterExecutionService(a.db)
	executionHandler := handler.NewEncounterExecutionHandler(executionService)

	router.HandleFunc("", executionHandler.Create).Methods("POST")
	router.HandleFunc("/get-all", executionHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", executionHandler.GetByID).Methods("GET")
	router.HandleFunc("/{id}", executionHandler.UpdateByID).Methods("PUT")
	router.HandleFunc("/{id}", executionHandler.DeleteByID).Methods("DELETE")
}

func (a *App) loadSocialEncounterRoutes(router *mux.Router) {
	socialEncounterService := service.NewSocialEncounterService(a.db)
	socialEncounterHandler := handler.NewSocialEncounterHandler(socialEncounterService)

	router.HandleFunc("", socialEncounterHandler.Create).Methods("POST")
	router.HandleFunc("/get-all", socialEncounterHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", socialEncounterHandler.GetByID).Methods("GET")
	router.HandleFunc("/{id}", socialEncounterHandler.UpdateByID).Methods("PUT")
	router.HandleFunc("/{id}", socialEncounterHandler.DeleteByID).Methods("DELETE")
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
