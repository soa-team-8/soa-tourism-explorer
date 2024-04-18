package handler

import (
	"context"
	"encoding/json"
	"followers/model"
	"followers/service"
	"net/http"
)

type contextKey string

const (
	ContextKeyUser contextKey = "user"
)

type SocialProfileHandler struct {
	service *service.SocialProfileService
}

func NewSocialProfileHandler(service *service.SocialProfileService) *SocialProfileHandler {
	return &SocialProfileHandler{service}
}

func (handler *SocialProfileHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ContextKeyUser).(*model.User)
	err := handler.service.CreateUser(user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (handler *SocialProfileHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func (handler *SocialProfileHandler) MiddlewareUserDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := &model.User{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyUser, user)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
