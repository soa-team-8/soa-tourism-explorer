package service

import (
	"encounters/dto"
	"encounters/repo"
	"fmt"
)

type EncounterRequestService struct {
	EncounterRequestRepo repo.EncounterRequestRepository
	EncounterRepo        repo.EncounterRepository
}

func NewEncounterRequestService(encounterRequestRepo repo.EncounterRequestRepository, encounterRepo repo.EncounterRepository) *EncounterRequestService {
	return &EncounterRequestService{
		EncounterRequestRepo: encounterRequestRepo,
		EncounterRepo:        encounterRepo,
	}
}

func (service *EncounterRequestService) Create(encounterReqDto dto.EncounterRequestDto) (dto.EncounterRequestDto, error) {
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

func (service *EncounterRequestService) Accept(id int) (dto.EncounterRequestDto, error) {
	request, err := service.EncounterRequestRepo.FindByID(id)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("request with ID %d not found", id)
	}

	request.Accept()

	updatedRequest, err := service.EncounterRequestRepo.Update(*request)

	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("request cannot be updated: %v", err)
	}

	encounterToPublish, err := service.EncounterRepo.FindByID(updatedRequest.EncounterId)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter not found: %v", err)
	}

	encounterToPublish.Publish()

	_, err = service.EncounterRepo.Update(*encounterToPublish)

	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter cannot be published: %v", err)
	}

	acceptedRequestDto := dto.ToDtoReq(*updatedRequest)
	return acceptedRequestDto, nil
}

func (service *EncounterRequestService) Reject(id int) (dto.EncounterRequestDto, error) {

	request, err := service.EncounterRequestRepo.FindByID(id)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("request with ID %d not found", id)
	}

	request.Reject()

	updatedRequest, err := service.EncounterRequestRepo.Update(*request)

	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("request cannot be updated: %v", err)
	}

	rejectedRequestDto := dto.ToDtoReq(*updatedRequest)
	return rejectedRequestDto, nil
}

func (service *EncounterRequestService) GetByID(id int) (dto.EncounterRequestDto, error) {
	encounterRequest, err := service.EncounterRequestRepo.FindByID(id)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter request not found: %v", err)
	}

	encounterRequestDto := dto.ToDtoReq(*encounterRequest)
	return encounterRequestDto, nil
}

func (service *EncounterRequestService) Update(encounterReqDto dto.EncounterRequestDto) (dto.EncounterRequestDto, error) {
	encounterReq := encounterReqDto.ToReqModel()
	updatedRequest, err := service.EncounterRequestRepo.Update(encounterReq)
	if err != nil {
		return dto.EncounterRequestDto{}, fmt.Errorf("encounter request cannot be updated: %v", err)
	}

	updatedRequestDto := dto.ToDtoReq(*updatedRequest)
	return updatedRequestDto, nil
}

func (service *EncounterRequestService) DeleteByID(id int) error {
	err := service.EncounterRequestRepo.DeleteByID(id)
	if err != nil {
		return fmt.Errorf("encounter request cannot be deleted: %v", err)
	}
	return nil
}
