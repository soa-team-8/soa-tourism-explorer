package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
	"tours/application"
	"tours/model"
	"tours/proto/checkpoints"
	"tours/proto/equipments"
	"tours/proto/tours"
	"tours/repository"
	"tours/service"
)

func main() {
	app := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	postgresClient := app.Db
	//EncounterService = createEncounterService(mongoClient)

	lis, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	tours.RegisterToursServiceServer(grpcServer, &Server{TourService: createTourService(postgresClient)})
	checkpoints.RegisterCheckpointsServiceServer(grpcServer, &Server{CheckpointService: createCheckpointService(postgresClient)})
	equipments.RegisterEquipmentsServiceServer(grpcServer, &Server{EquipmentService: createEquipmentService(postgresClient)})

	reflection.Register(grpcServer)
	grpcServer.Serve(lis)

	defer cancel()

	err = app.Start(ctx)
	if err != nil {
		fmt.Println("Failed to start app:", err)
	}
}

func createTourService(postgresClient *gorm.DB) *service.TourService {
	tourService := &service.TourService{
		TourRepository: &repository.TourRepository{
			DB: postgresClient,
		},
		EquipmentRepository: &repository.EquipmentRepository{
			DB: postgresClient,
		},
	}
	return tourService
}

func createCheckpointService(postgresClient *gorm.DB) *service.CheckpointService {
	checkpointService := &service.CheckpointService{
		CheckpointRepository: &repository.CheckpointRepository{
			DB: postgresClient,
		},
	}
	return checkpointService
}

func createEquipmentService(postgresClient *gorm.DB) *service.EquipmentService {
	equipmentService := &service.EquipmentService{
		EquipmentRepository: &repository.EquipmentRepository{
			DB: postgresClient,
		},
	}
	return equipmentService
}

type Server struct {
	tours.UnimplementedToursServiceServer
	checkpoints.UnimplementedCheckpointsServiceServer
	equipments.UnimplementedEquipmentsServiceServer
	TourService       *service.TourService
	CheckpointService *service.CheckpointService
	EquipmentService  *service.EquipmentService
}

// Metoda za kreiranje ture
func (s *Server) CreateTour(ctx context.Context, in *tours.TourDto) (*tours.ActionResponse, error) {
	// Prvo mapirajte osnovne informacije
	var dem model.DemandignessLevel
	if in.GetDemandingnessLevel() == "Easy" {
		dem = model.Easy
	} else if in.GetDemandingnessLevel() == "Medium" {
		dem = model.Medium
	} else if in.GetDemandingnessLevel() == "Hard" {
		dem = model.Hard
	}

	var status model.Status
	if in.GetStatus() == "Draft" {
		status = model.Draft
	} else if in.GetStatus() == "Published" {
		status = model.Published
	} else if in.GetStatus() == "Archived" {
		status = model.Archived
	}

	tour := &model.Tour{
		ID:                uint64(in.GetId()),
		AuthorID:          uint64(in.GetAuthorId()),
		Name:              in.GetName(),
		Description:       in.GetDescription(),
		DemandignessLevel: dem,
		Status:            status,
		Price:             in.GetPrice(),
		Tags:              in.GetTags(),
		Closed:            in.GetClosed(),
	}

	// Mapirajte opremu (equipment)
	var equipment []model.Equipment
	for _, equipmentDto := range in.GetEquipment() {
		equipment = append(equipment, model.Equipment{
			ID:          uint64(equipmentDto.GetId()),
			Name:        equipmentDto.GetName(),
			Description: equipmentDto.GetDescription(),
		})
	}
	tour.Equipment = equipment

	// Mapirajte checkpointe
	var checkpoints []model.Checkpoint
	for _, checkpointDto := range in.GetCheckpoints() {
		checkpoint := model.Checkpoint{
			ID:                    uint64(checkpointDto.GetId()),
			TourID:                uint64(checkpointDto.GetTourId()),
			AuthorID:              uint64(checkpointDto.GetAuthorId()),
			Longitude:             checkpointDto.GetLongitude(),
			Latitude:              checkpointDto.GetLatitude(),
			Name:                  checkpointDto.GetName(),
			Description:           checkpointDto.GetDescription(),
			Pictures:              checkpointDto.GetPictures(),
			RequiredTimeInSeconds: checkpointDto.GetRequiredTimeInSeconds(),
			EncounterID:           uint64(checkpointDto.GetEncounterId()),
			IsSecretPrerequisite:  checkpointDto.GetIsSecretPrerequisite(),
		}
		// Mapiranje CheckpointSecretDto
		checkpointSecretDto := checkpointDto.GetCheckpointSecret()
		checkpointSecret := model.CheckpointSecret{
			Description: checkpointSecretDto.GetDescription(),
			Pictures:    checkpointSecretDto.GetPictures(),
		}
		checkpoint.CheckpointSecret = checkpointSecret

		checkpoints = append(checkpoints, checkpoint)
	}
	tour.Checkpoints = checkpoints

	// Mapirajte published ture
	var publishedTours []model.PublishedTour
	for _, publishedTourDto := range in.GetPublishedTours() {
		publishedTour := model.PublishedTour{
			PublishingDate: time.Unix(publishedTourDto.GetPublishingDate().GetSeconds(), int64(publishedTourDto.GetPublishingDate().GetNanos())),
		}
		publishedTours = append(publishedTours, publishedTour)
	}
	tour.PublishedTours = publishedTours

	// Mapirajte archived ture
	var archivedTours []model.ArchivedTour
	for _, archivedTourDto := range in.GetArchivedTours() {
		archivedTour := model.ArchivedTour{
			ArchivingDate: time.Unix(archivedTourDto.GetArchivingDate().GetSeconds(), int64(archivedTourDto.GetArchivingDate().GetNanos())),
		}
		archivedTours = append(archivedTours, archivedTour)
	}
	tour.ArchivedTours = archivedTours

	s.TourService.Create(*tour)

	// Vratite odgovor sa statusom uspeha
	return &tours.ActionResponse{Succes: true}, nil
}

func (s *Server) GetAllTours(ctx context.Context, in *tours.GetAllToursRequest) (*tours.ListTourDtoResponse, error) {
	var toursReal []model.Tour
	toursReal, _ = s.TourService.GetAll()

	// Kreiramo praznu listu DTO tura
	var tourDtos []*tours.TourDto

	// Mapiramo modele tura u DTO oblike
	for _, tour := range toursReal {
		tourDto := &tours.TourDto{
			Id:                 int32(tour.ID),
			Name:               tour.Name,
			Description:        tour.Description,
			DemandingnessLevel: tour.DemandignessLevel.String(), // Konvertujemo enum u string
			Price:              tour.Price,
			AuthorId:           int32(tour.AuthorID),
			Status:             tour.Status.String(), // Konvertujemo enum u string
			Tags:               tour.Tags,
			Closed:             tour.Closed,
			// Ostale atribute tura možete dodati ovde
		}

		// Mapiranje opreme (equipment)
		for _, equipment := range tour.Equipment {
			equipmentDto := &tours.EquipmentDto{
				Id:          int32(equipment.ID),
				Name:        equipment.Name,
				Description: equipment.Description,
			}
			tourDto.Equipment = append(tourDto.Equipment, equipmentDto)
		}

		// Mapiranje checkpointa
		for _, checkpoint := range tour.Checkpoints {
			checkpointDto := &tours.CheckpointDto{
				Id:                    int64(checkpoint.ID),
				TourId:                int64(checkpoint.TourID),
				AuthorId:              int64(checkpoint.AuthorID),
				Longitude:             checkpoint.Longitude,
				Latitude:              checkpoint.Latitude,
				Name:                  checkpoint.Name,
				Description:           checkpoint.Description,
				Pictures:              checkpoint.Pictures,
				RequiredTimeInSeconds: checkpoint.RequiredTimeInSeconds,
				EncounterId:           int64(checkpoint.EncounterID),
				IsSecretPrerequisite:  checkpoint.IsSecretPrerequisite,
			}

			checkpointDto.CheckpointSecret = &tours.CheckpointSecretDto{
				Description: checkpoint.CheckpointSecret.Description,
				Pictures:    checkpoint.CheckpointSecret.Pictures,
			}

			tourDto.Checkpoints = append(tourDto.Checkpoints, checkpointDto)
		}

		// Dodavanje DTO tura u listu
		tourDtos = append(tourDtos, tourDto)
	}

	// Vraćamo odgovor sa listom DTO tura
	return &tours.ListTourDtoResponse{Tours: tourDtos}, nil
}

func (s *Server) AddEquipment(ctx context.Context, in *tours.AddEquipmentRequest) (*tours.ActionResponse, error) {
	_ = s.TourService.AddEquipmentToTour(uint64(in.TourId), uint64(in.EquipmentId))
	return &tours.ActionResponse{Succes: true}, nil
}

func (s *Server) CreateCheckpoint(ctx context.Context, in *checkpoints.CreateCheckpointRequest) (*checkpoints.CheckpointDto2, error) {
	// Prvo mapirajte osnovne informacije
	checkpoint := model.Checkpoint{
		TourID:                uint64(in.GetCheckpoint().GetTourId()),
		AuthorID:              uint64(in.GetCheckpoint().GetAuthorId()),
		Longitude:             in.GetCheckpoint().GetLongitude(),
		Latitude:              in.GetCheckpoint().GetLatitude(),
		Name:                  in.GetCheckpoint().GetName(),
		Description:           in.GetCheckpoint().GetDescription(),
		Pictures:              in.GetCheckpoint().GetPictures(),
		RequiredTimeInSeconds: in.GetCheckpoint().GetRequiredTimeInSeconds(),
		EncounterID:           uint64(in.GetCheckpoint().GetEncounterId()),
		IsSecretPrerequisite:  in.GetCheckpoint().GetIsSecretPrerequisite(),
	}

	// Mapiranje CheckpointSecretDto
	checkpointSecretDto := in.GetCheckpoint().GetCheckpointSecret()
	checkpointSecret := model.CheckpointSecret{
		Description: checkpointSecretDto.GetDescription(),
		Pictures:    checkpointSecretDto.GetPictures(),
	}
	checkpoint.CheckpointSecret = checkpointSecret

	// Pozovite funkciju koja će dodati checkpoint u bazu podataka
	_, err := s.CheckpointService.Create(checkpoint)
	if err != nil {
		// U slučaju greške, vratite odgovarajući odgovor sa greškom
		return nil, err
	}

	// Vratite odgovor sa kreiranim checkpoint-om
	return in.Checkpoint, nil
}

func (s *Server) CreateCheckpointSecret(ctx context.Context, in *checkpoints.CreateCheckpointSecretRequest) (*checkpoints.CheckpointDto2, error) {
	secret := model.CheckpointSecret{
		Description: in.CheckpointSecret.GetDescription(),
		Pictures:    in.CheckpointSecret.GetPictures(),
	}

	// Poziv metode za kreiranje ili ažuriranje tajnog checkpoint-a u servisu
	err := s.CheckpointService.CreateOrUpdateCheckpointSecret(uint64(in.CheckpointId), secret)
	if err != nil {
		// Obrada greške ako je potrebno
		return nil, err
	}
	updatedCheckpoint, err := s.CheckpointService.GetByID(uint64(in.CheckpointId))
	if err != nil {
		// Obrada greške ako je potrebno
		return nil, err
	}

	// Mapiranje ažuriranog checkpoint-a na odgovarajući DTO
	updatedCheckpointDto := checkpoints.CheckpointDto2{
		Id:                    int64(updatedCheckpoint.ID),
		TourId:                int64(updatedCheckpoint.TourID),
		AuthorId:              int64(updatedCheckpoint.AuthorID),
		Longitude:             updatedCheckpoint.Longitude,
		Latitude:              updatedCheckpoint.Latitude,
		Name:                  updatedCheckpoint.Name,
		Description:           updatedCheckpoint.Description,
		Pictures:              updatedCheckpoint.Pictures,
		RequiredTimeInSeconds: updatedCheckpoint.RequiredTimeInSeconds,
		EncounterId:           int64(updatedCheckpoint.EncounterID),
		IsSecretPrerequisite:  updatedCheckpoint.IsSecretPrerequisite,
	}

	// Mapiranje CheckpointSecret na odgovarajući DTO polje
	updatedCheckpointDto.CheckpointSecret = &checkpoints.CheckpointSecretDto2{
		Description: secret.Description,
		Pictures:    secret.Pictures,
	}

	// Vraćanje DTO-a ažuriranog checkpoint-a kao odgovor
	return &updatedCheckpointDto, nil
}

func (s *Server) GetAllEquipment(ctx context.Context, in *equipments.GetAllEquipmentRequest) (*equipments.EquipmentListDto, error) {
	// Dobijamo stvarnu listu opreme iz servisa
	equipmentsReal, err := s.EquipmentService.GetAll()
	if err != nil {
		// Obrada greške ako je potrebno
		return nil, err
	}

	// Kreiramo praznu listu DTO opreme
	var eqDtos []*equipments.EquipmentDto3

	// Mapiranje stvarne opreme na DTO opremu
	for _, equipment := range equipmentsReal {
		eqDto := &equipments.EquipmentDto3{
			Id:          int32(equipment.ID),
			Name:        equipment.Name,
			Description: equipment.Description,
		}
		eqDtos = append(eqDtos, eqDto)
	}

	// Kreiranje EquipmentListDto sa mapiranom opremom
	equipmentListDto := &equipments.EquipmentListDto{
		Items: eqDtos,
	}

	// Vraćamo listu DTO opreme kao odgovor
	return equipmentListDto, nil
}
