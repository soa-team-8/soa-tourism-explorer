package main

import (
	"context"
	"encounters/application"
	"encounters/dto"
	"encounters/model"
	"encounters/proto/encounters"
	"encounters/repo/mongoDB"
	"encounters/service"
	"fmt"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

//var EncounterService *service.EncounterService

func main() {
	app := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	mongoClient := app.MongoClient
	//EncounterService = createEncounterService(mongoClient)

	lis, err := net.Listen("tcp", "localhost:3030")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	encounters.RegisterEncounterServiceServer(grpcServer, Server{EncounterService: createEncounterService(mongoClient)})

	reflection.Register(grpcServer)
	grpcServer.Serve(lis)

	defer cancel()

	err = app.Start(ctx)
	if err != nil {
		fmt.Println("Failed to start app:", err)
	}

}

type Server struct {
	encounters.UnimplementedEncounterServiceServer
	EncounterService *service.EncounterService
}

func (s Server) CreateEncounter(ctx context.Context, request *encounters.CreateEncounterRequest) (*encounters.EncounterResponse, error) {
	encounterDto := request.GetEncounter()

	images := pq.StringArray{encounterDto.Image}
	requiredPeople := int(encounterDto.RequiredPeople)
	activeTouristsIDsBigSlice := toBigIntSlice(encounterDto.ActiveTouristsIds)

	_, err := s.EncounterService.CreateByAuthor(dto.EncounterDto{
		AuthorID:          uint64(encounterDto.AuthorId),
		ID:                uint64(encounterDto.Id),
		Name:              encounterDto.Name,
		Description:       encounterDto.Description,
		XP:                encounterDto.Xp,
		Status:            encounterDto.Status,
		Type:              encounterDto.Type,
		Longitude:         float64(encounterDto.Longitude),
		Latitude:          float64(encounterDto.Latitude),
		LocationLatitude:  &encounterDto.LocationLatitude,
		LocationLongitude: &encounterDto.LocationLongitude,
		Image:             images,
		Range:             &encounterDto.Range,
		RequiredPeople:    &requiredPeople,
		ActiveTouristsIDs: activeTouristsIDsBigSlice,
	})

	if err != nil {
		return nil, err
	}

	response := &encounters.EncounterResponse{
		Encounter: encounterDto,
	}

	return response, nil

}

func toBigIntSlice(ids []int32) *model.BigIntSlice {
	var bigIntSlice model.BigIntSlice
	for _, id := range ids {
		bigIntSlice = append(bigIntSlice, uint64(id))
	}
	return &bigIntSlice
}

func createEncounterService(mongoClient *mongo.Client) *service.EncounterService {
	encounterRepository := mongoDB.NewEncounterRepository(mongoClient)
	encounterRequestRepository := mongoDB.NewEncounterRequestRepository(mongoClient)
	socialEncounterRepository := mongoDB.NewSocialEncounterRepository(mongoClient)
	locationEncounterRepository := mongoDB.NewHiddenLocationRepository(mongoClient)
	return service.NewEncounterService(encounterRepository, encounterRequestRepository, socialEncounterRepository, locationEncounterRepository)
}
