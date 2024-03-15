package application

import (
	"fmt"
	"net/http"
	"time"
	"tours/handler"
	"tours/repository"
	"tours/service"

	"github.com/gorilla/mux"
)

func (a *App) loadRoutes() {
	router := mux.NewRouter()
	router.Use(loggerMiddleware)

	toursRouter := router.PathPrefix("/tours").Subrouter()
	a.loadTourRoutes(toursRouter)

	equipmentRouter := router.PathPrefix("/equipment").Subrouter()
	a.loadEquipmentRoutes(equipmentRouter)

	checkpointRouter := router.PathPrefix("/checkpoints").Subrouter()
	a.loadCheckpointRoutes(checkpointRouter)

	router.HandleFunc("/images/{imageName}", a.serveImage).Methods("GET")

	a.router = router
}

func (a *App) loadTourRoutes(router *mux.Router) {
	tourService := &service.TourService{
		TourRepository: &repository.TourRepository{
			DB: a.db,
		},
		EquipmentRepository: &repository.EquipmentRepository{
			DB: a.db,
		},
	}

	tourHandler := &handler.TourHandler{
		TourService: tourService,
	}

	router.HandleFunc("", tourHandler.Create).Methods("POST")
	router.HandleFunc("", tourHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", tourHandler.Update).Methods("PUT")
	router.HandleFunc("/{id}", tourHandler.Delete).Methods("DELETE")
	router.HandleFunc("/{id}", tourHandler.GetByID).Methods("GET")
	router.HandleFunc("/{tourID}/{equipmentID}/add", tourHandler.AddEquipmentToTour).Methods("PUT")
	router.HandleFunc("/{tourID}/{equipmentID}/remove", tourHandler.RemoveEquipmentFromTour).Methods("PUT")
	router.HandleFunc("/{authorID}/by-author", tourHandler.GetToursByAuthor).Methods("GET")
	router.HandleFunc("/{id}/publish", tourHandler.Publish).Methods("PUT")
	router.HandleFunc("/{id}/archive", tourHandler.Archive).Methods("PUT")
}

func (a *App) loadEquipmentRoutes(router *mux.Router) {
	equipmentService := &service.EquipmentService{
		EquipmentRepository: &repository.EquipmentRepository{
			DB: a.db,
		},
	}

	equipmentHandler := &handler.EquipmentHandler{
		EquipmentService: equipmentService,
	}

	router.HandleFunc("", equipmentHandler.Create).Methods("POST")
	router.HandleFunc("", equipmentHandler.GetAll).Methods("GET")
	router.HandleFunc("/paged", equipmentHandler.GetAllPaged).Methods("GET")
	router.HandleFunc("/{id}", equipmentHandler.Update).Methods("PUT")
	router.HandleFunc("/{id}", equipmentHandler.Delete).Methods("DELETE")
	router.HandleFunc("/{id}", equipmentHandler.GetByID).Methods("GET")
	router.HandleFunc("/{id}/get-available", equipmentHandler.GetAvailableEquipment).Methods("POST")
}

func (a *App) loadCheckpointRoutes(router *mux.Router) {
	checkpointService := &service.CheckpointService{
		CheckpointRepository: &repository.CheckpointRepository{
			DB: a.db,
		},
	}

	checkpointHandler := &handler.CheckpointHandler{
		CheckpointService: checkpointService,
	}

	router.HandleFunc("", checkpointHandler.Create).Methods("POST")
	router.HandleFunc("", checkpointHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", checkpointHandler.Update).Methods("PUT")
	router.HandleFunc("/{id}", checkpointHandler.Delete).Methods("DELETE")
	router.HandleFunc("/{id}", checkpointHandler.GetByID).Methods("GET")
	router.HandleFunc("/{id}/tour", checkpointHandler.GetAllByTourID).Methods("GET")
	router.HandleFunc("/{id}/checkpoint-secret", checkpointHandler.CreateOrUpdateCheckpointSecret).Methods("PUT")
	router.HandleFunc("/setEnc/{checkpointId}/{encId}/{isSecretPrerequisite}", checkpointHandler.SetCheckpointEncounter).Methods("PUT")
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
