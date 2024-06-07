package main

import (
	"context"
	"encounters/application"
	"encounters/dto"
	"encounters/model"
	"encounters/proto/encounter_requests"
	"encounters/proto/encounters"
	"encounters/proto/tourist_encounters"
	"encounters/repo/mongoDB"
	"encounters/service"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	netsock "net"
	"os"
	"os/signal"
	"time"
)

// var loggingFile, err = os.OpenFile("logging/var/encounter-api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
var logger = log.New(os.Stdout, "[encounter-api] ", log.LstdFlags)

var (
	cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_cpu_usage_percent",
		Help: "Current CPU usage percentage",
	})
	memUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_mem_usage_percent",
		Help: "Current memory usage percentage",
	})
	diskUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_disk_usage_percent",
		Help: "Current disk usage percentage",
	})
	netSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "host_network_bytes_sent_total",
		Help: "Total bytes sent over the network",
	})
	netRecv = promauto.NewCounter(prometheus.CounterOpts{
		Name: "host_network_bytes_received_total",
		Help: "Total bytes received over the network",
	})
)

func collectMetrics() {
	for {
		// Collect CPU usage
		cpuPercent, err := cpu.Percent(time.Second, false)
		if err == nil && len(cpuPercent) > 0 {
			cpuUsage.Set(cpuPercent[0])
		}

		// Collect memory usage
		virtualMem, err := mem.VirtualMemory()
		if err == nil {
			memUsage.Set(virtualMem.UsedPercent)
		}

		// Collect disk usage
		diskInfo, err := disk.Usage("/")
		if err == nil {
			diskUsage.Set(diskInfo.UsedPercent)
		}

		// Collect network usage
		netIO, err := net.IOCounters(false)
		if err == nil && len(netIO) > 0 {
			netSent.Add(float64(netIO[0].BytesSent))
			netRecv.Add(float64(netIO[0].BytesRecv))
		}

		time.Sleep(10 * time.Second) // Adjust the collection interval as needed
	}
}

func main() {
	app := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/metrics", promhttp.Handler())
	go collectMetrics()

	mongoClient := app.MongoClient
	//EncounterService = createEncounterService(mongoClient)

	lis, err := netsock.Listen("tcp", "localhost:3031")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	server := Server{EncounterService: createEncounterService(mongoClient), EncounterRequestService: createEncounterRequestService(mongoClient)}

	encounters.RegisterEncounterServiceServer(grpcServer, server)
	tourist_encounters.RegisterTouristEncounterServiceServer(grpcServer, server)
	encounter_requests.RegisterEncounterRequestServiceServer(grpcServer, server)

	reflection.Register(grpcServer)
	grpcServer.Serve(lis)

	logger.Println("API")

	defer cancel()

	err = app.Start(ctx)
	if err != nil {
		fmt.Println("Failed to start app:", err)
	}

}

type Server struct {
	encounters.UnimplementedEncounterServiceServer
	tourist_encounters.UnimplementedTouristEncounterServiceServer
	encounter_requests.UnimplementedEncounterRequestServiceServer

	EncounterService        *service.EncounterService
	EncounterRequestService *service.EncounterRequestService
}

func (s Server) CreateEncounter(ctx context.Context, request *encounters.CreateEncounterRequest) (*encounters.EncounterResponse, error) {
	encounterDto := request.GetEncounter()

	images := pq.StringArray{encounterDto.Image}
	requiredPeople := int(encounterDto.RequiredPeople)
	activeTouristsIDsBigSlice := toBigIntSlice(encounterDto.ActiveTouristsIds)

	logger.Printf("Received CreateEncounter request: %+v", encounterDto)

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
		logger.Printf("Error creating encounter: %v", err)
		return nil, err
	}

	response := &encounters.EncounterResponse{
		Encounter: encounterDto,
	}

	logger.Printf("Successfully created encounter: %+v", response)

	return response, nil
}

func (s Server) DeleteEncounter(ctx context.Context, request *encounters.DeleteEncounterRequest) (*encounters.DeleteEncounterResponse, error) {
	if err := s.EncounterService.DeleteByID(uint64(request.GetId())); err != nil {
		return nil, err
	}

	return &encounters.DeleteEncounterResponse{Success: true}, nil
}

func (s Server) UpdateEncounter(ctx context.Context, request *encounters.UpdateEncounterRequest) (*encounters.EncounterResponse, error) {
	encounterDto := request.GetEncounter()

	images := pq.StringArray{encounterDto.Image}
	requiredPeople := int(encounterDto.RequiredPeople)
	activeTouristsIDsBigSlice := toBigIntSlice(encounterDto.ActiveTouristsIds)

	_, err := s.EncounterService.Update(dto.EncounterDto{
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

func (s Server) GetEncounter(ctx context.Context, request *encounters.GetEncounterRequest) (*encounters.EncounterResponse, error) {
	foundedEncounter, err := s.EncounterService.GetByID(uint64(request.GetId()))

	if err != nil {
		return nil, err
	}

	var activeTouristsIDs []int32 = nil
	var locationLongitude float64 = 0
	var locationLatitude float64 = 0
	var inRange float64 = 0
	var image = ""
	var requiredPeople int32 = 0

	if foundedEncounter.LocationLongitude != nil {
		locationLongitude = *foundedEncounter.LocationLongitude
	}

	if foundedEncounter.LocationLatitude != nil {
		locationLatitude = *foundedEncounter.LocationLatitude
	}

	if foundedEncounter.Range != nil {
		inRange = *foundedEncounter.Range
	}

	if foundedEncounter.Image != nil {
		image = foundedEncounter.Image[0]
	}

	if foundedEncounter.RequiredPeople != nil {
		requiredPeople = int32(*foundedEncounter.RequiredPeople)
	}

	if foundedEncounter.ActiveTouristsIDs != nil {
		for _, id := range *foundedEncounter.ActiveTouristsIDs {
			activeTouristsIDs = append(activeTouristsIDs, int32(id))
		}
	}

	encounter := &encounters.Encounter{
		AuthorId:          int64(foundedEncounter.AuthorID),
		Id:                int64(foundedEncounter.ID),
		Name:              foundedEncounter.Name,
		Description:       foundedEncounter.Description,
		Xp:                foundedEncounter.XP,
		Status:            foundedEncounter.Status,
		Type:              foundedEncounter.Type,
		Longitude:         float32(foundedEncounter.Longitude),
		Latitude:          float32(foundedEncounter.Latitude),
		LocationLongitude: locationLongitude,
		LocationLatitude:  locationLatitude,
		Image:             image,
		Range:             inRange,
		RequiredPeople:    requiredPeople,
		ActiveTouristsIds: activeTouristsIDs,
	}

	response := &encounters.EncounterResponse{
		Encounter: encounter,
	}

	return response, nil
}

func (s Server) TouristCreateEncounter(ctx context.Context, request *tourist_encounters.TouristCreateEncounterRequest) (*tourist_encounters.TouristCreateEncounterResponse, error) {
	encounterDto := request.GetEncounter()
	level := 11
	var userId uint64 = 1

	images := pq.StringArray{encounterDto.Image}
	requiredPeople := int(encounterDto.RequiredPeople)
	activeTouristsIDsBigSlice := toBigIntSlice(encounterDto.ActiveTouristsIds)

	_, err := s.EncounterService.CreateByTourist(dto.EncounterDto{
		AuthorID:          uint64(encounterDto.AuthorId),
		ID:                uint64(encounterDto.Id),
		Name:              encounterDto.Name,
		Description:       encounterDto.Description,
		XP:                encounterDto.XP,
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
	}, level, userId)

	if err != nil {
		return nil, err
	}

	response := &tourist_encounters.TouristCreateEncounterResponse{
		Encounter: encounterDto,
	}

	return response, nil
}

func (s Server) TouristGetAllEncounters(ctx context.Context, request *tourist_encounters.TouristGetAllEncountersRequest) (*tourist_encounters.TouristGetAllEncountersResponse, error) {
	encounters, err := s.EncounterService.GetAll()
	if err != nil {
		return nil, err
	}

	var responseEncounters []*tourist_encounters.TouristEncounter

	for _, foundedEncounter := range encounters {
		var activeTouristsIDs []int32
		var locationLongitude float64
		var locationLatitude float64
		var inRange float64
		var image string
		var requiredPeople int32

		if foundedEncounter.LocationLongitude != nil {
			locationLongitude = *foundedEncounter.LocationLongitude
		}

		if foundedEncounter.LocationLatitude != nil {
			locationLatitude = *foundedEncounter.LocationLatitude
		}

		if foundedEncounter.Range != nil {
			inRange = *foundedEncounter.Range
		}

		if foundedEncounter.Image != nil && len(foundedEncounter.Image) > 0 {
			image = foundedEncounter.Image[0]
		}

		if foundedEncounter.RequiredPeople != nil {
			requiredPeople = int32(*foundedEncounter.RequiredPeople)
		}

		if foundedEncounter.ActiveTouristsIDs != nil {
			for _, id := range *foundedEncounter.ActiveTouristsIDs {
				activeTouristsIDs = append(activeTouristsIDs, int32(id))
			}
		}

		encounter := &tourist_encounters.TouristEncounter{
			AuthorId:          int64(foundedEncounter.AuthorID),
			Id:                int64(foundedEncounter.ID),
			Name:              foundedEncounter.Name,
			Description:       foundedEncounter.Description,
			XP:                int32(foundedEncounter.XP),
			Status:            foundedEncounter.Status,
			Type:              foundedEncounter.Type,
			Longitude:         float32(foundedEncounter.Longitude),
			Latitude:          float32(foundedEncounter.Latitude),
			LocationLongitude: locationLongitude,
			LocationLatitude:  locationLatitude,
			Image:             image,
			Range:             inRange,
			RequiredPeople:    requiredPeople,
			ActiveTouristsIds: activeTouristsIDs,
		}

		responseEncounters = append(responseEncounters, encounter)
	}

	return &tourist_encounters.TouristGetAllEncountersResponse{Encounters: responseEncounters}, nil
}

func (s Server) CreateEncounterRequest(ctx context.Context, request *encounter_requests.CreateEncounterRequestDto) (*encounter_requests.EncounterRequestResponseDto, error) {
	encounterRequest := request.EncounterRequest

	encounterRequestDto := dto.EncounterRequestDto{
		ID:          uint64(encounterRequest.Id),
		EncounterId: uint64(encounterRequest.EncounterId),
		TouristId:   uint64(encounterRequest.TouristId),
		Status:      encounterRequest.Status,
	}

	_, err := s.EncounterRequestService.Create(encounterRequestDto)

	if err != nil {
		return nil, err
	}

	response := &encounter_requests.EncounterRequestResponseDto{
		EncounterRequest: encounterRequest,
	}

	return response, nil
}

func (s Server) GetAllEncounterRequests(ctx context.Context, request *encounter_requests.GetAllEncounterRequestsRequest) (*encounter_requests.GetAllEncounterRequestsResponse, error) {
	encounterRequests, err := s.EncounterRequestService.GetAll()

	if err != nil {
		return nil, err
	}

	var encounterRequestsResp []*encounter_requests.EncounterRequestDto

	for _, encounterReq := range encounterRequests {
		encReq := encounter_requests.EncounterRequestDto{
			Id:          int32(encounterReq.ID),
			EncounterId: int64(encounterReq.EncounterId),
			TouristId:   int64(encounterReq.TouristId),
			Status:      encounterReq.Status,
		}

		encounterRequestsResp = append(encounterRequestsResp, &encReq)
	}

	return &encounter_requests.GetAllEncounterRequestsResponse{EncounterRequests: encounterRequestsResp}, nil
}

func (s Server) AcceptEncounterRequest(ctx context.Context, request *encounter_requests.AcceptEncounterRequestDto) (*encounter_requests.EncounterRequestResponseDto, error) {
	acceptedReq, err := s.EncounterRequestService.Accept(int(request.Id))
	if err != nil {
		return nil, err
	}

	encReq := encounter_requests.EncounterRequestDto{
		Id:          int32(acceptedReq.ID),
		EncounterId: int64(acceptedReq.EncounterId),
		TouristId:   int64(acceptedReq.TouristId),
		Status:      acceptedReq.Status,
	}

	response := encounter_requests.EncounterRequestResponseDto{
		EncounterRequest: &encReq,
	}

	return &response, nil
}

func (s Server) RejectEncounterRequest(ctx context.Context, request *encounter_requests.RejectEncounterRequestDto) (*encounter_requests.EncounterRequestResponseDto, error) {
	acceptedReq, err := s.EncounterRequestService.Reject(int(request.Id))
	if err != nil {
		return nil, err
	}

	encReq := encounter_requests.EncounterRequestDto{
		Id:          int32(acceptedReq.ID),
		EncounterId: int64(acceptedReq.EncounterId),
		TouristId:   int64(acceptedReq.TouristId),
		Status:      acceptedReq.Status,
	}

	response := encounter_requests.EncounterRequestResponseDto{
		EncounterRequest: &encReq,
	}

	return &response, nil
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

func createEncounterRequestService(mongoClient *mongo.Client) *service.EncounterRequestService {
	encounterRequestRepository := mongoDB.NewEncounterRequestRepository(mongoClient)
	encounterRepository := mongoDB.NewEncounterRepository(mongoClient)
	return service.NewEncounterRequestService(encounterRequestRepository, encounterRepository)
}
