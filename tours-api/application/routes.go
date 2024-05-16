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

	tourRatingRouter := router.PathPrefix("/tour-ratings").Subrouter()
	a.loadTourRatingRoutes(tourRatingRouter)

	tourExecutionRouter := router.PathPrefix("/tour-executions").Subrouter()
	a.loadTourExecutionRoutes(tourExecutionRouter)

	publishedTourRouter := router.PathPrefix("/published-tours").Subrouter()
	a.loadPublishedToursRoutes(publishedTourRouter)

	reportedIssueRouter := router.PathPrefix("/reported-issues").Subrouter()
	a.loadReportedIssuesRoutes(reportedIssueRouter)

	router.HandleFunc("/images/{imageName}", a.serveImage).Methods("GET")

	a.router = router
}

func (a *App) loadTourRoutes(router *mux.Router) {
	tourService := &service.TourService{
		TourRepository: &repository.TourRepository{
			DB: a.Db,
		},
		EquipmentRepository: &repository.EquipmentRepository{
			DB: a.Db,
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
			DB: a.Db,
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
			DB: a.Db,
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
	router.HandleFunc("/get-encounter-ids/{tourId}", checkpointHandler.GetEncounterIDsByTour).Methods("GET")
}

func (a *App) loadTourRatingRoutes(router *mux.Router) {
	tourRatingService := &service.TourRatingService{
		TourRatingRepository: &repository.TourRatingRepository{
			DB: a.Db,
		},
		TourExecutionRepository: &repository.TourExecutionRepository{
			DB: a.Db,
		},
	}

	tourRatingHandler := &handler.TourRatingHandler{
		TourRatingService: tourRatingService,
	}

	router.HandleFunc("", tourRatingHandler.Create).Methods("POST")
	router.HandleFunc("", tourRatingHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", tourRatingHandler.Update).Methods("PUT")
	router.HandleFunc("/{id}", tourRatingHandler.Delete).Methods("DELETE")
	router.HandleFunc("/{id}", tourRatingHandler.GetByID).Methods("GET")
}

func (a *App) loadTourExecutionRoutes(router *mux.Router) {
	tourExecutionService := &service.TourExecutionService{
		TourExecutionRepository: &repository.TourExecutionRepository{
			DB: a.Db,
		},
	}

	tourExecutionHandler := &handler.TourExecutionHandler{
		TourExecutionService: tourExecutionService,
	}

	router.HandleFunc("/{id}", tourExecutionHandler.CheckPosition).Methods("PUT")
	router.HandleFunc("/{userID}/{executionID}", tourExecutionHandler.Abandon).Methods("PUT")
	router.HandleFunc("/{userID}/{tourID}", tourExecutionHandler.Create).Methods("POST")
	router.HandleFunc("/{userID}/{tourID}", tourExecutionHandler.GetByIDs).Methods("GET")
}

func (a *App) loadReportedIssuesRoutes(router *mux.Router) {
	reportedIssueService := &service.ReportedIssueService{
		ReportedIssueRepository: &repository.ReportedIssueRepository{
			DB: a.Db,
		},
		TourRepository: &repository.TourRepository{
			DB: a.Db,
		},
	}

	reportedIssueHandler := &handler.ReportedIssueHandler{
		ReportedIssueService: reportedIssueService,
	}

	router.HandleFunc("", reportedIssueHandler.GetAll).Methods("GET")
	router.HandleFunc("/{id}/author", reportedIssueHandler.GetAllByAuthor).Methods("GET")
	router.HandleFunc("/{id}/tourist", reportedIssueHandler.GetAllByTourist).Methods("GET")
	router.HandleFunc("/{cat}/{desc}/{prior}/{tourID}/{userID}", reportedIssueHandler.Create).Methods("POST")
	router.HandleFunc("/{id}/comment", reportedIssueHandler.AddComment).Methods("POST")
	router.HandleFunc("/{id}/deadline", reportedIssueHandler.AddDeadline).Methods("PUT")
	router.HandleFunc("/{id}/penalize", reportedIssueHandler.PenalizeAuthor).Methods("PUT")
	router.HandleFunc("/{id}/close", reportedIssueHandler.Close).Methods("PUT")
	router.HandleFunc("/{id}/resolve", reportedIssueHandler.Resolve).Methods("PUT")

}

func (a *App) loadPublishedToursRoutes(router *mux.Router) {
	tourService := &service.TourService{
		TourRepository: &repository.TourRepository{
			DB: a.Db,
		},
		TourRatingRepository: &repository.TourRatingRepository{
			DB: a.Db,
		},
	}

	publishedToursHandler := &handler.PublishedToursHandler{
		TourService: tourService,
	}

	router.HandleFunc("", publishedToursHandler.GetPublishedTours).Methods("GET")
	router.HandleFunc("/{id}", publishedToursHandler.GetPublishedTour).Methods("GET")
	router.HandleFunc("/{id}/rating", publishedToursHandler.GetAverageRating).Methods("GET")
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
