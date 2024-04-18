package main

import (
	"context"
	"errors"
	"followers/application"
	handlers "followers/handler"
	repositories "followers/repository"
	"followers/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[followers-api] ", log.LstdFlags)

	repo, err := repositories.NewSocialProfileRepository()
	if err != nil {
		logger.Fatal(err)
	}
	defer repo.CloseDriverConnection(timeoutContext)

	socialProfileService := service.NewSocialProfileService(repo)
	SocialProfileHandler := handlers.NewSocialProfileHandler(socialProfileService)

	router := mux.NewRouter()
	router.Use(SocialProfileHandler.MiddlewareContentTypeSet)
	application.SetUpRoutes(router, SocialProfileHandler)

	server := configureServer(router)
	startServer(server, logger)
	shutdownServer(server, logger)
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
