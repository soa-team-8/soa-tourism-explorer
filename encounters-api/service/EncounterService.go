package service

import (
	"encounters/dto"
	"encounters/model"
	repo "encounters/repo"
	"fmt"
)

type EncounterService struct {
	EncounterRepo *repo.EncounterRepository
}

func (service *EncounterService) Create(dto dto.EncounterDto) error {
	encounter := dto.ToModel()

	err := service.EncounterRepo.Save(encounter)
	if err != nil {
		return fmt.Errorf("encounter cannot be created: %v", err)
	}
	return nil
}

func (service *EncounterService) GetByID(id uint64) (*model.Encounter, error) {
	encounter, err := service.EncounterRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Encounter with id %d not found", id))
	}
	return encounter, nil
}

func (service *EncounterService) GetAll() ([]model.Encounter, error) {
	encounters, err := service.EncounterRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Encounters not found"))
	}
	return encounters, nil
}

func (service *EncounterService) DeleteByID(id uint64) error {
	err := service.EncounterRepo.DeleteByID(id)
	if err != nil {
		return fmt.Errorf("encounter cannot be deleted: %v", err)
	}
	return nil
}

func (service *EncounterService) Update(dto dto.EncounterDto) error {
	encounter := dto.ToModel()
	err := service.EncounterRepo.Update(encounter)
	if err != nil {
		return fmt.Errorf("encounter cannot be updated: %v", err)
	}
	return nil
}
