package service

import (
	"encounters/dto"
	"encounters/repo"
	"fmt"
	"gorm.io/gorm"
)

type SocialEncounterService struct {
	SocialEncounterRepo *repo.SocialEncounterRepository
}

func NewSocialEncounterService(db *gorm.DB) *SocialEncounterService {
	return &SocialEncounterService{
		SocialEncounterRepo: &repo.SocialEncounterRepository{
			Db: db,
		},
	}
}

func (service *SocialEncounterService) Create(encounterDto dto.EncounterDto) (dto.EncounterDto, error) {
	encounter := encounterDto.ToSocialModel()

	savedEncounter, err := service.SocialEncounterRepo.Save(encounter)
	if err != nil {
		return dto.EncounterDto{}, fmt.Errorf("encounter cannot be created: %v", err)
	}

	savedEncounterDto := dto.ToSocialDto(savedEncounter)

	return savedEncounterDto, nil
}

func (service *SocialEncounterService) GetByID(id uint64) (*dto.EncounterDto, error) {
	encounter, err := service.SocialEncounterRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("encounter with ID %d not found", id)
	}

	encounterDto := dto.ToSocialDto(*encounter)
	return &encounterDto, nil
}

func (service *SocialEncounterService) GetAll() ([]dto.EncounterDto, error) {
	encounters, err := service.SocialEncounterRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Encounters not found"))
	}

	encounterDtos := dto.ToSocialDtoList(encounters)
	return encounterDtos, nil
}

func (service *SocialEncounterService) DeleteByID(id uint64) error {
	err := service.SocialEncounterRepo.DeleteById(id)
	if err != nil {
		return fmt.Errorf("encounter cannot be deleted: %v", err)
	}
	return nil
}

func (service *SocialEncounterService) Update(encounterDto dto.EncounterDto) (dto.EncounterDto, error) {
	encounter := encounterDto.ToSocialModel()

	updatedEncounter, err := service.SocialEncounterRepo.Update(encounter)
	if err != nil {
		return dto.EncounterDto{}, fmt.Errorf("encounter cannot be updated: %v", err)
	}

	updatedEncounterDto := dto.ToSocialDto(updatedEncounter)

	return updatedEncounterDto, nil
}
