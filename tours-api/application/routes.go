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

	a.router = router
}

func (a *App) loadTourRoutes(router *mux.Router) {
	tourService := &service.TourService{
		TourRepository: &repository.TourRepository{
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
