package main

import (
	"context"
	"errors"
	"followers-api/application"
	handlers "followers-api/handler"
	social_profile "followers-api/proto/social-profile"
	repositories "followers-api/repository"
	"followers-api/service"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	logger := log.New(os.Stdout, "[followers-api] ", log.LstdFlags)

	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	repo, err := repositories.NewSocialProfileRepository()
	if err != nil {
		logger.Fatal(err)
	}

	socialProfileService := service.NewSocialProfileService(repo)
	SocialProfileHandler := handlers.NewSocialProfileHandler(socialProfileService)

	socialProfilegRPCServer := Server{SocialProfileService: socialProfileService}

	social_profile.RegisterSocialProfileServiceServer(grpcServer, socialProfilegRPCServer)

	reflection.Register(grpcServer)
	grpcServer.Serve(lis)

	defer cancel()

	router := mux.NewRouter()
	router.Use(SocialProfileHandler.MiddlewareContentTypeSet)
	application.SetUpRoutes(router, SocialProfileHandler)

	defer repo.CloseDriverConnection(timeoutContext)

	server := configureServer(router)
	startServer(server, logger)
	shutdownServer(server, logger)
}

type Server struct {
	social_profile.UnimplementedSocialProfileServiceServer
	SocialProfileService *service.SocialProfileService
}

func configureServer(router http.Handler) *http.Server {
	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	return &http.Server{
		Addr:         ":9090",
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}

func startServer(server *http.Server, logger *log.Logger) {
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Could not start server: %v", err)
		}
	}()
}

func shutdownServer(server *http.Server, logger *log.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)

	sig := <-sigCh
	logger.Printf("Received signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Graceful shutdown failed: %v", err)
	}

	logger.Println("Server stopped")
}

func (s Server) Follow(ctx context.Context, request *social_profile.FollowRequest) (*social_profile.FollowResponse, error) {
	followerID := request.FollowerId
	followedID := request.FollowedId

	updatedProfile, err := s.SocialProfileService.Follow(uint64(followerID), uint64(followedID))

	if err != nil {
		return nil, err
	}

	var followersList []*social_profile.UserDto
	var followedList []*social_profile.UserDto
	var followableList []*social_profile.UserDto

	for _, follower := range updatedProfile.Followers {
		followerDto := social_profile.UserDto{
			Id:       int32(follower.ID),
			Username: follower.Username,
		}

		followersList = append(followersList, &followerDto)
	}

	for _, followed := range updatedProfile.Followed {
		followerDto := social_profile.UserDto{
			Id:       int32(followed.ID),
			Username: followed.Username,
		}

		followedList = append(followedList, &followerDto)
	}

	for _, followable := range updatedProfile.Followable {
		followerDto := social_profile.UserDto{
			Id:       int32(followable.ID),
			Username: followable.Username,
		}

		followableList = append(followableList, &followerDto)
	}

	response := &social_profile.FollowResponse{
		SocialProfile: &social_profile.SocialProfileDto{
			Id:         int32(updatedProfile.UserId),
			Followed:   followersList,
			Followers:  followedList,
			Followable: followableList,
		},
	}

	return response, nil
}

func (s Server) UnFollow(ctx context.Context, request *social_profile.UnFollowRequest) (*social_profile.UnFollowResponse, error) {
	followerID := request.FollowerId
	followedID := request.UnfollowedId

	updatedProfile, err := s.SocialProfileService.Unfollow(uint64(followerID), uint64(followedID))

	if err != nil {
		return nil, err
	}

	var followersList []*social_profile.UserDto
	var followedList []*social_profile.UserDto
	var followableList []*social_profile.UserDto

	for _, follower := range updatedProfile.Followers {
		followerDto := social_profile.UserDto{
			Id:       int32(follower.ID),
			Username: follower.Username,
		}

		followersList = append(followersList, &followerDto)
	}

	for _, followed := range updatedProfile.Followed {
		followerDto := social_profile.UserDto{
			Id:       int32(followed.ID),
			Username: followed.Username,
		}

		followedList = append(followedList, &followerDto)
	}

	for _, followable := range updatedProfile.Followable {
		followerDto := social_profile.UserDto{
			Id:       int32(followable.ID),
			Username: followable.Username,
		}

		followableList = append(followableList, &followerDto)
	}

	response := &social_profile.UnFollowResponse{
		SocialProfile: &social_profile.SocialProfileDto{
			Id:         int32(updatedProfile.UserId),
			Followed:   followersList,
			Followers:  followedList,
			Followable: followableList,
		},
	}

	return response, nil
}

func (s Server) GetSocialProfile(ctx context.Context, request *social_profile.GetSocialProfileRequest) (*social_profile.GetSocialProfileResponse, error) {
	userId := request.UserId
	profile, err := s.SocialProfileService.GetProfile(uint64(userId))

	if err != nil {
		return nil, err
	}

	var followersList []*social_profile.UserDto
	var followedList []*social_profile.UserDto
	var followableList []*social_profile.UserDto

	for _, follower := range profile.Followers {
		followerDto := social_profile.UserDto{
			Id:       int32(follower.ID),
			Username: follower.Username,
		}

		followersList = append(followersList, &followerDto)
	}

	for _, followed := range profile.Followed {
		followerDto := social_profile.UserDto{
			Id:       int32(followed.ID),
			Username: followed.Username,
		}

		followedList = append(followedList, &followerDto)
	}

	for _, followable := range profile.Followable {
		followerDto := social_profile.UserDto{
			Id:       int32(followable.ID),
			Username: followable.Username,
		}

		followableList = append(followableList, &followerDto)
	}

	response := &social_profile.GetSocialProfileResponse{
		SocialProfile: &social_profile.SocialProfileDto{
			Id:         int32(profile.UserId),
			Followed:   followersList,
			Followers:  followedList,
			Followable: followableList,
		},
	}

	return response, nil

}

func (s Server) GetRecommendations(ctx context.Context, request *social_profile.GetRecommendationsRequest) (*social_profile.GetRecommendationsResponse, error) {
	userId := request.GetUserId()

	recommendations, err := s.SocialProfileService.GetRecommendations(uint64(userId))
	if err != nil {
		return nil, err
	}

	var recommendationsList []*social_profile.UserDto

	for _, recommendation := range recommendations {
		userDto := social_profile.UserDto{
			Id:       int32(recommendation.ID),
			Username: recommendation.Username,
		}
		recommendationsList = append(recommendationsList, &userDto)

	}

	response := &social_profile.GetRecommendationsResponse{
		Recommendations: recommendationsList,
	}

	return response, nil
}
