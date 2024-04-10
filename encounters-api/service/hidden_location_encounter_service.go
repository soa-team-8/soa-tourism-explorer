package service

import (
	"encounters/dto"
	"encounters/repo/postgreSQL"
	"fmt"
	"gorm.io/gorm"
)

type HiddenLocationEncounterService struct {
	HiddenLocationEncounterRepo *postgreSQL.HiddenLocationRepository
}

func NewHiddenLocationEncounterService(db *gorm.DB) *HiddenLocationEncounterService {
	return &HiddenLocationEncounterService{
		HiddenLocationEncounterRepo: &postgreSQL.HiddenLocationRepository{
			Db: db,
		},
	}
}

func (service *HiddenLocationEncounterService) Create(encounterDto dto.EncounterDto) (dto.EncounterDto, error) {
	encounter := encounterDto.ToHiddenLocationModel()

	savedEncounter, err := service.HiddenLocationEncounterRepo.Save(encounter)
	if err != nil {
		return dto.EncounterDto{}, fmt.Errorf("encounter cannot be created: %v", err)
	}

	savedEncounterDto := dto.ToHiddenLocationDto(savedEncounter)

	return savedEncounterDto, nil
}

func (service *HiddenLocationEncounterService) GetByID(id uint64) (*dto.EncounterDto, error) {
	encounter, err := service.HiddenLocationEncounterRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("encounter with ID %d not found", id)
	}

	encounterDto := dto.ToHiddenLocationDto(*encounter)
	return &encounterDto, nil
}

func (service *HiddenLocationEncounterService) GetAll() ([]dto.EncounterDto, error) {
	encounters, err := service.HiddenLocationEncounterRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Encounters not found"))
	}

	encounterDtos := dto.ToHiddenLocationDtoList(encounters)
	return encounterDtos, nil
}

func (service *HiddenLocationEncounterService) DeleteByID(id uint64) error {
	err := service.HiddenLocationEncounterRepo.DeleteById(id)
	if err != nil {
		return fmt.Errorf("encounter cannot be deleted: %v", err)
	}
	return nil
}

func (service *HiddenLocationEncounterService) Update(encounterDto dto.EncounterDto) (dto.EncounterDto, error) {
	encounter := encounterDto.ToHiddenLocationModel()

	updatedEncounter, err := service.HiddenLocationEncounterRepo.Update(encounter)
	if err != nil {
		return dto.EncounterDto{}, fmt.Errorf("encounter cannot be updated: %v", err)
	}

	updatedEncounterDto := dto.ToHiddenLocationDto(updatedEncounter)

	return updatedEncounterDto, nil
}
