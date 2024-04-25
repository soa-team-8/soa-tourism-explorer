package service

import (
	"encounters/dto"
	"encounters/model"
	"encounters/repo"
	"fmt"
)

type EncounterService struct {
	EncounterRepo        repo.EncounterRepository
	EncounterRequestRepo repo.EncounterRequestRepository
	SocialEncounterRepo  repo.SocialEncounterRepository
	HiddenEncounterRepo  repo.HiddenLocationRepository
}

func NewEncounterService(encounterRepo repo.EncounterRepository, encounterRequestRepo repo.EncounterRequestRepository,
	socialEncounterRepo repo.SocialEncounterRepository, hiddenEncounterRepo repo.HiddenLocationRepository) *EncounterService {
	return &EncounterService{
		EncounterRequestRepo: encounterRequestRepo,
		EncounterRepo:        encounterRepo,
		SocialEncounterRepo:  socialEncounterRepo,
		HiddenEncounterRepo:  hiddenEncounterRepo,
	}
}

func (service *EncounterService) Create(encounterDto dto.EncounterDto) (dto.EncounterDto, error) {
	encounter := encounterDto.ToModel()

	savedEncounter, err := service.EncounterRepo.Save(encounter)
	if err != nil {
		return dto.EncounterDto{}, fmt.Errorf("encounter cannot be created: %v", err)
	}

	savedEncounterDto := dto.ToDto(savedEncounter)

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

func (service *EncounterService) CreateByTourist(encounterDto dto.EncounterDto, level int, userId uint64) (dto.EncounterDto, error) {
	var savedEncId uint64
	if level >= 10 {
		if encounterDto.Type == "Location" {
			var hiddenLocationEncounter = encounterDto.ToHiddenLocationModel()
			savedEncounter, err := service.HiddenEncounterRepo.Save(hiddenLocationEncounter)
			savedEncId = savedEncounter.ID
			if err != nil {
				return dto.EncounterDto{}, fmt.Errorf("hidden location encounter cannot be created: %v", err)
			}
		} else if encounterDto.Type == "Social" {
			var socialEncounter = encounterDto.ToSocialModel()
			savedEncounter, err := service.SocialEncounterRepo.Save(socialEncounter)
			savedEncId = savedEncounter.ID
			if err != nil {
				return dto.EncounterDto{}, fmt.Errorf("social encounter cannot be created: %v", err)
			}
		} else {
			var encounter = encounterDto.ToModel()
			savedEncounter, err := service.EncounterRepo.Save(encounter)
			savedEncId = savedEncounter.ID
			if err != nil {
				return dto.EncounterDto{}, fmt.Errorf("encounter cannot be created: %v", err)
			}
		}
		encounterDto.ID = savedEncId
		encounterReqDto := dto.EncounterRequestDto{TouristId: userId, EncounterId: savedEncId, Status: "OnHold"}
		_, err := service.EncounterRequestRepo.Save(encounterReqDto.ToReqModel())
		if err != nil {
			return dto.EncounterDto{}, err
		}
		return encounterDto, err
	} else {
		return encounterDto, fmt.Errorf("the tourist is not at level 10 or higher")
	}
}

func (service *EncounterService) CreateByAuthor(encounterDto dto.EncounterDto) (dto.EncounterDto, error) {
	var savedEncounterDto dto.EncounterDto

	if encounterDto.Type == "Location" {
		var hiddenLocationEncounter = encounterDto.ToHiddenLocationModel()
		savedEncounter, err := service.HiddenEncounterRepo.Save(hiddenLocationEncounter)
		if err != nil {
			return savedEncounterDto, fmt.Errorf("hidden location encounter cannot be created: %v", err)
		}
		savedEncounterDto = dto.ToHiddenLocationDto(savedEncounter)
	} else if encounterDto.Type == "Social" {
		var socialEncounter = encounterDto.ToSocialModel()
		savedEncounter, err := service.SocialEncounterRepo.Save(socialEncounter)
		if err != nil {
			return savedEncounterDto, fmt.Errorf("social encounter cannot be created: %v", err)
		}
		savedEncounterDto = dto.ToSocialDto(savedEncounter)
	} else {
		var encounter = encounterDto.ToModel()
		savedEncounter, err := service.EncounterRepo.Save(encounter)
		if err != nil {
			return savedEncounterDto, fmt.Errorf("encounter cannot be created: %v", err)
		}
		savedEncounterDto = dto.ToDto(savedEncounter)
	}

	return savedEncounterDto, nil
}

func (service *EncounterService) AddEncounter(execution model.EncounterExecution) (model.EncounterExecution, error) {
	newEncounter, err := service.EncounterRepo.FindByID(execution.EncounterID)
	if err != nil {
		return model.EncounterExecution{}, fmt.Errorf("encounter with ID %d not found", execution.EncounterID)
	}

	execution.Encounter = *newEncounter

	return execution, nil
}

func (service *EncounterService) AddEncounters(executions []model.EncounterExecution) ([]model.EncounterExecution, error) {
	var addedExecutions []model.EncounterExecution

	for _, execution := range executions {
		newExecution, err := service.AddEncounter(execution)
		if err != nil {
			return nil, fmt.Errorf("error adding encounter: %v", err)
		}
		addedExecutions = append(addedExecutions, newExecution)
	}

	return addedExecutions, nil
}
