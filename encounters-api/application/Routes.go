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
