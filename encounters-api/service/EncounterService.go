package service

import (
	"encounters/dto"
	"encounters/repo"
	"fmt"
)

type EncounterService struct {
	EncounterRepo *repo.EncounterRepository
}

func (service *EncounterService) Create(encounterDto dto.EncounterDto) (dto.EncounterDto, error) {
	encounter := encounterDto.ToModel()

	savedEncounter, err := service.EncounterRepo.Save(encounter)
	if err != nil {
		return dto.EncounterDto{}, fmt.Errorf("encounter cannot be created: %v", err)
	}

	savedEncounterDto := dto.ToDto(savedEncounter)

	fmt.Println(savedEncounter.ID)
	return savedEncounterDto, nil
}

func (service *EncounterService) GetByID(id uint64) (*dto.EncounterDto, error) {
	encounter, err := service.EncounterRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("encounter with ID %d not found", id)
	}

	encounterDto := dto.ToDto(*encounter)
	return &encounterDto, nil
}

func (service *EncounterService) GetAll() ([]dto.EncounterDto, error) {
	encounters, err := service.EncounterRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Encounters not found"))
	}

	encounterDtos := dto.ToDtoList(encounters)
	return encounterDtos, nil
}

func (service *EncounterService) DeleteByID(id uint64) error {
	err := service.EncounterRepo.DeleteByID(id)
	if err != nil {
		return fmt.Errorf("encounter cannot be deleted: %v", err)
	}
	return nil
}

func (service *EncounterService) Update(encounterDto dto.EncounterDto) (dto.EncounterDto, error) {
	encounter := encounterDto.ToModel()

	updatedEncounter, err := service.EncounterRepo.Update(encounter)
	if err != nil {
		return dto.EncounterDto{}, fmt.Errorf("encounter cannot be updated: %v", err)
	}

	updatedEncounterDto := dto.ToDto(updatedEncounter)

	return updatedEncounterDto, nil
}
