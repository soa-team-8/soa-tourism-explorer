package application

import (
	"encounters/handler"
	"encounters/repo"
	"encounters/repo/mongoDB"
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

	hiddenLocationEncounterRouter := router.PathPrefix("/encounters/hiddenLoc").Subrouter()
	a.loadHiddenLocationEncounterRoutes(hiddenLocationEncounterRouter)

	a.router = router
}

func (a *App) loadEncounterRoutes(router *mux.Router) {
	encounterService := service.NewEncounterService(a.postgresDB)
	encounterHandler := handler.NewEncounterHandler(encounterService)

	router.HandleFunc("", encounterHandler.Create).Methods("POST")
	router.HandleFunc("/get-all", encounterHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", encounterHandler.GetByID).Methods("GET")
	router.HandleFunc("/{id}", encounterHandler.UpdateByID).Methods("PUT")
	router.HandleFunc("/{id}", encounterHandler.DeleteByID).Methods("DELETE")
	router.HandleFunc("/tourist/{level}/{userId}", encounterHandler.CreateTouristEncounter).Methods("POST")
	router.HandleFunc("/author", encounterHandler.CreateAuthorEncounter).Methods("POST")
}

func (a *App) loadEncounterRequestRoutes(router *mux.Router) {
	encounterRequestRepository := mongoDB.New(a.mongoClient)
	encounterRepository := repo.New(a.postgresDB)
	encounterRequestService := service.NewEncounterRequestService(encounterRequestRepository, *encounterRepository)
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
	executionService := service.NewEncounterExecutionService(a.postgresDB)
	encounterService := service.NewEncounterService(a.postgresDB)
	executionHandler := handler.NewEncounterExecutionHandler(executionService, encounterService)

	router.HandleFunc("/{touristId}", executionHandler.Create).Methods("POST")
	router.HandleFunc("/get-all", executionHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", executionHandler.GetByID).Methods("GET")
	router.HandleFunc("/{touristId}/{id}", executionHandler.UpdateByID).Methods("PUT")
	router.HandleFunc("/{touristId}/{id}", executionHandler.DeleteByID).Methods("DELETE")

	//Complex routes
	router.HandleFunc("/activate/{touristId}/{encounterId}", executionHandler.Activate).Methods("PUT")
	router.HandleFunc("/complete/{touristId}/{executionId}", executionHandler.Complete).Methods("PUT")

	router.HandleFunc("/get-by-tour/{touristId}", executionHandler.GetByTour).Methods("GET")
	router.HandleFunc("/get-active-by-tour/{touristId}", executionHandler.GetActiveByTour).Methods("GET")
	router.HandleFunc("/get-all-by-tourist/{touristId}", executionHandler.GetAllByTourist).Methods("GET")
	router.HandleFunc("/get-completed-by-tourist/{touristId}", executionHandler.GetAllCompletedByTourist).Methods("GET")

	router.HandleFunc("/social-encounter/check-range/{encounterId}/{tourId}/{touristId}", executionHandler.CheckPosition).Methods("GET")
	router.HandleFunc("/location-encounter/check-range/{encounterId}/{tourId}/{touristId}", executionHandler.CheckPositionLocationEncounter).
		Methods("GET")
}

func (a *App) loadSocialEncounterRoutes(router *mux.Router) {
	socialEncounterService := service.NewSocialEncounterService(a.postgresDB)
	socialEncounterHandler := handler.NewSocialEncounterHandler(socialEncounterService)

	router.HandleFunc("", socialEncounterHandler.Create).Methods("POST")
	router.HandleFunc("/get-all", socialEncounterHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", socialEncounterHandler.GetByID).Methods("GET")
	router.HandleFunc("/{id}", socialEncounterHandler.UpdateByID).Methods("PUT")
	router.HandleFunc("/{id}", socialEncounterHandler.DeleteByID).Methods("DELETE")
}

func (a *App) loadHiddenLocationEncounterRoutes(router *mux.Router) {
	hiddenLocationEncounterService := service.NewHiddenLocationEncounterService(a.postgresDB)
	hiddenLocationEncounterHandler := handler.NewHiddenLocationEncounterHandler(hiddenLocationEncounterService)

	router.HandleFunc("", hiddenLocationEncounterHandler.Create).Methods("POST")
	router.HandleFunc("/get-all", hiddenLocationEncounterHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", hiddenLocationEncounterHandler.GetByID).Methods("GET")
	router.HandleFunc("/{id}", hiddenLocationEncounterHandler.UpdateByID).Methods("PUT")
	router.HandleFunc("/{id}", hiddenLocationEncounterHandler.DeleteByID).Methods("DELETE")
}

func (a *App) serveImage(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	imageName := vars["imageName"]

	http.ServeFile(resp, req, "wwwroot/images/"+imageName)
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
