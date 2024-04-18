package application

import (
	handlers "followers/handler"
	"github.com/gorilla/mux"
	"net/http"
)

func SetUpRoutes(router *mux.Router, handler *handlers.SocialProfileHandler) {
	socialProfileRouter := router.PathPrefix("/social-profile").Subrouter()

	socialProfileRouter.HandleFunc("/user", handler.CreateUser).Methods(http.MethodPost)
	socialProfileRouter.HandleFunc("/follow/{followerId}/{followedId}", handler.Follow).Methods(http.MethodPost)
	socialProfileRouter.Use(handler.MiddlewareUserDeserialization)
}
