package application

import (
	handlers "followers-api/handler"
	"github.com/gorilla/mux"
	"net/http"
)

func SetUpRoutes(router *mux.Router, handler *handlers.SocialProfileHandler) {
	socialProfileRouter := router.PathPrefix("/social-profile").Subrouter()

	socialProfileRouter.HandleFunc("/user", handler.CreateUser).Methods(http.MethodPost)
	socialProfileRouter.Use(handler.MiddlewareUserDeserialization)
	socialProfileRouter.HandleFunc("/{id}", handler.GetProfile).Methods(http.MethodGet)
	socialProfileRouter.HandleFunc("/recommendations/{id}", handler.GetRecommendations).Methods(http.MethodGet)
	socialProfileRouter.HandleFunc("/follow/{followerId}/{followedId}", handler.Follow).Methods(http.MethodPost)
	socialProfileRouter.HandleFunc("/unfollow/{followerId}/{followedId}", handler.Unfollow).Methods(http.MethodPost)
}
