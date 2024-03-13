package service

import (
	"encounters/dto"
	"encounters/repo"
	"fmt"
)

// EncounterRequestService je servis za rad sa zahtevima za susrete
type EncounterRequestService struct {
	EncounterRequestRepo *repo.EncounterRequestRepository
	EncounterService     *EncounterService
	EncounterRepo        *repo.EncounterRepository
}

// NewEncounterRequestService kreira novi EncounterRequestService
func NewEncounterRequestService(encounterRequestRepo *repo.EncounterRequestRepository) *EncounterRequestService {
	return &EncounterRequestService{EncounterRequestRepo: encounterRequestRepo}
}

// CreateEncounterRequest kreira novi zahtev za susret
func (service *EncounterRequestService) CreateEncounterRequest(encounterReqDto dto.EncounterRequestDto) (dto.EncounterRequestDto, error) {
	encounterReq := encounterReqDto.ToReqModel()
	newRequest, err := service.EncounterRequestRepo.Save(encounterReq)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter request cannot be created: %v", err)
	}

	newRequestDto := dto.ToDtoReq(newRequest)
	return newRequestDto, nil
}

func (service *EncounterRequestService) GetAll() ([]dto.EncounterRequestDto, error) {
	encounterRequests, err := service.EncounterRequestRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Encounters not found"))
	}

	encounterRequestsDtos := dto.ToDtoListReq(encounterRequests)
	return encounterRequestsDtos, nil
}

// AcceptEncounterRequest prihvata zahtev za susret sa datim ID-om
func (service *EncounterRequestService) AcceptEncounterRequest(id int) (dto.EncounterRequestDto, error) {
	acceptedRequest, err := service.EncounterRequestRepo.AcceptRequest(id)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter request cannot be accepted: %v", err)
	}
	encouterToPublishDto, err := service.EncounterService.GetByID(acceptedRequest.EncounterId)
	service.EncounterRepo.MakeEncounterPublished(encouterToPublishDto.ID)
	acceptedRequestDto := dto.ToDtoReq(*acceptedRequest)
	return acceptedRequestDto, nil
}

// RejectEncounterRequest odbija zahtev za susret sa datim ID-om
func (service *EncounterRequestService) RejectEncounterRequest(id int) (dto.EncounterRequestDto, error) {
	rejectedRequest, err := service.EncounterRequestRepo.RejectRequest(id)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter request cannot be rejected: %v", err)
	}

	rejectedRequestDto := dto.ToDtoReq(*rejectedRequest)
	return rejectedRequestDto, nil
}
