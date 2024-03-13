package service

import (
	"encounters/dto"
	"encounters/repo"
	"fmt"
)

// EncounterRequestService je servis za rad sa zahtevima za susrete
type EncounterRequestService struct {
	EncounterRequestRepo *repo.EncounterRequestRepository
	EncounterRepo        *repo.EncounterRepository
}

// NewEncounterRequestService kreira novi EncounterRequestService
func NewEncounterRequestService(encounterRequestRepo *repo.EncounterRequestRepository, encounterRepo *repo.EncounterRepository) *EncounterRequestService {
	return &EncounterRequestService{
		EncounterRequestRepo: encounterRequestRepo,
		EncounterRepo:        encounterRepo,
	}
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

	encounterToPublish, err := service.EncounterRepo.FindByID(acceptedRequest.EncounterId)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter not found: %v", err)
	}

	_, err = service.EncounterRepo.MakeEncounterPublished(encounterToPublish.ID)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter cannot be published: %v", err)
	}

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

// GetEncounterRequestByID pronalazi zahtev za susret sa datim ID-om
func (service *EncounterRequestService) GetEncounterRequestByID(id int) (dto.EncounterRequestDto, error) {
	encounterRequest, err := service.EncounterRequestRepo.FindByID(id)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter request not found: %v", err)
	}

	encounterRequestDto := dto.ToDtoReq(*encounterRequest)
	return encounterRequestDto, nil
}

// UpdateEncounterRequest ažurira postojeći zahtev za susret
func (service *EncounterRequestService) UpdateEncounterRequest(encounterReqDto dto.EncounterRequestDto) (dto.EncounterRequestDto, error) {
	encounterReq := encounterReqDto.ToReqModel()
	updatedRequest, err := service.EncounterRequestRepo.Update(encounterReq)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter request cannot be updated: %v", err)
	}

	updatedRequestDto := dto.ToDtoReq(*updatedRequest)
	return updatedRequestDto, nil
}

// DeleteEncounterRequestByID briše zahtev za susret sa datim ID-om
func (service *EncounterRequestService) DeleteEncounterRequestByID(id int) error {
	err := service.EncounterRequestRepo.DeleteByID(id)
	if err != nil {
		return fmt.Errorf("encounter request cannot be deleted: %v", err)
	}
	return nil
}
