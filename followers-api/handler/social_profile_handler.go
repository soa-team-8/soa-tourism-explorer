package handler

import (
	"context"
	"encoding/json"
	"followers/model"
	"followers/service"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func (handler *SocialProfileHandler) Follow(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	followerID, err := strconv.ParseUint(params["followerId"], 10, 64)
	if err != nil {
		http.Error(rw, "Invalid follower ID", http.StatusBadRequest)
		return
	}
	followedID, err := strconv.ParseUint(params["followedId"], 10, 64)
	if err != nil {
		http.Error(rw, "Invalid followed ID", http.StatusBadRequest)
		return
	}

	updatedProfile, err := handler.service.Follow(followerID, followedID)
	if err != nil {
		http.Error(rw, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(updatedProfile)
}

func (handler *SocialProfileHandler) Unfollow(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	followerID, err := strconv.ParseUint(params["followerId"], 10, 64)
	if err != nil {
		http.Error(rw, "Invalid follower ID", http.StatusBadRequest)
		return
	}
	followedID, err := strconv.ParseUint(params["followedId"], 10, 64)
	if err != nil {
		http.Error(rw, "Invalid followed ID", http.StatusBadRequest)
		return
	}

	updatedProfile, err := handler.service.Unfollow(followerID, followedID)
	if err != nil {
		http.Error(rw, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(updatedProfile)
}

func (handler *SocialProfileHandler) GetProfile(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}

	profile, err := handler.service.GetProfile(id)
	if err != nil {
		http.Error(rw, "Failed to get user", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(rw).Encode(profile)
	if err != nil {
		return
	}
}

func (handler *SocialProfileHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func (handler *SocialProfileHandler) MiddlewareUserDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.ContentLength > 0 {
			user := &model.User{}
			err := json.NewDecoder(r.Body).Decode(user)
			if err != nil {
				http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
				return
			}
			ctx := context.WithValue(r.Context(), ContextKeyUser, user)
			next.ServeHTTP(rw, r.WithContext(ctx))
		} else {
			next.ServeHTTP(rw, r)
		}
	})
}
