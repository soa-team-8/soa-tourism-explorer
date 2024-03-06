package application

import (
	"context"
	"fmt"
	"net/http"

	"time"

	"encounters/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	router http.Handler
	db     *gorm.DB
	config Config
}

func New(config Config) *App {
	connectionURL := config.PostgresAddress
	db, err := gorm.Open(postgres.Open(connectionURL), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Perform migrations
	err = db.AutoMigrate(&model.Encounter{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	app := &App{
		db:     db,
		config: config,
	}

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
