package application

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"

	"time"

	"encounters/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	router      http.Handler
	postgresDB  *gorm.DB
	mongoClient *mongo.Client
	config      Config
}

func New(config Config) *App {
	app := &App{
		config: config,
	}

	if err := app.setupPostgres(); err != nil {
		fmt.Println("Failed to setup PostgreSQL:", err)
		return nil
	}

	if err := app.setupMongoDB(); err != nil {
		fmt.Println("Failed to setup MongoDB:", err)
		return nil
	}

	app.loadRoutes()

	return app
}

func (a *App) setupPostgres() error {
	postgresConnectionURL := a.config.PostgresAddress
	db, err := gorm.Open(postgres.Open(postgresConnectionURL), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := performMigrations(db); err != nil {
		return err
	}

	a.postgresDB = db
	return nil
}

func (a *App) setupMongoDB() error {
	mongoURI := a.config.MongoDBAddress
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mongoClient.Connect(ctx)
	if err != nil {
		return err
	}

	if err := createEncountersDatabase(ctx, mongoClient); err != nil {
		return err
	}

	a.mongoClient = mongoClient
	return nil
}

func performMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Encounter{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.SocialEncounter{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.HiddenLocationEncounter{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.EncounterRequest{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.EncounterExecution{}); err != nil {
		return err
	}

	return nil
}

func createEncountersDatabase(ctx context.Context, client *mongo.Client) error {
	err := client.Database("encounters").CreateCollection(ctx, "encounterExecutions")
	err = client.Database("encounters").CreateCollection(ctx, "encounterRequests")
	err = client.Database("encounters").CreateCollection(ctx, "encounters")
	return err
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
