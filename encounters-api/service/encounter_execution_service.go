package service

import (
	"encounters/model"
	"encounters/repo"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type EncounterExecutionService struct {
	ExecutionRepo         *repo.EncounterExecutionRepository
	EncounterRepo         *repo.EncounterRepository
	SocialEncounterRepo   *repo.SocialEncounterRepository
	LocationEncounterRepo *repo.HiddenLocationRepository
}

func NewEncounterExecutionService(db *gorm.DB) *EncounterExecutionService {
	return &EncounterExecutionService{
		ExecutionRepo:         repo.NewEncounterExecutionRepository(db),
		EncounterRepo:         repo.New(db),
		SocialEncounterRepo:   repo.NewSocialEncounterRepository(db),
		LocationEncounterRepo: repo.NewHiddenLocationRepository(db),
	}
}

func (service *EncounterExecutionService) Create(execution model.EncounterExecution, touristID uint64) (model.EncounterExecution, error) {
	if execution.TouristID != touristID {
		return model.EncounterExecution{}, fmt.Errorf("encounter touristID does not match provided touristID")
	}

	savedExecution, err := service.ExecutionRepo.Save(execution)

	if err != nil {
		return model.EncounterExecution{}, fmt.Errorf("execution cannot be created: %v", err)
	}

	return savedExecution, nil
}

func (service *EncounterExecutionService) GetByID(id uint64) (*model.EncounterExecution, error) {
	execution, err := service.ExecutionRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("execution with ID %d not found", id)
	}

	return execution, nil
}

func (service *EncounterExecutionService) GetAll() ([]model.EncounterExecution, error) {
	executions, err := service.ExecutionRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Executions not found"))
	}

	return executions, nil
}

func (service *EncounterExecutionService) DeleteByID(id uint64, touristID uint64) error {
	// Check permission
	if err := service.checkPermission(id, touristID); err != nil {
		return err
	}

	// Delete the execution
	err := service.ExecutionRepo.DeleteByID(id)
	if err != nil {
		return fmt.Errorf("execution cannot be deleted: %v", err)
	}

	return nil
}

func (service *EncounterExecutionService) Update(execution model.EncounterExecution, touristID uint64) (model.EncounterExecution, error) {
	if err := service.checkPermission(execution.ID, touristID); err != nil {
		return model.EncounterExecution{}, err
	}

	updatedExecution, err := service.ExecutionRepo.Update(execution)
	if err != nil {
		return model.EncounterExecution{}, fmt.Errorf("execution cannot be updated: %v", err)
	}

	return updatedExecution, nil
}

func (service *EncounterExecutionService) GetAllByTourist(touristID uint64) ([]model.EncounterExecution, error) {
	executions, err := service.ExecutionRepo.FindAllByTourist(touristID)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Executions not found"))
	}

	return executions, nil
}

func (service *EncounterExecutionService) GetAllActiveByTourist(touristID uint64) ([]model.EncounterExecution, error) {
	executions, err := service.ExecutionRepo.FindAllByTourist(touristID)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Executions not found"))
	}

	return executions, nil
}

func (service *EncounterExecutionService) GetAllCompletedByTourist(touristID uint64) ([]model.EncounterExecution, error) {
	executions, err := service.ExecutionRepo.FindAllCompletedByTourist(touristID)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Executions not found"))
	}

	return executions, nil
}

func (service *EncounterExecutionService) GetByEncounter(encounterId uint64) (*model.EncounterExecution, error) {
	execution, err := service.ExecutionRepo.FindByEncounter(encounterId)
	if err != nil {
		return nil, fmt.Errorf("execution with ID %d not found", encounterId)
	}

	return execution, nil
}

func (service *EncounterExecutionService) GetAllByEncounter(encounterID uint64) ([]model.EncounterExecution, error) {
	executions, err := service.ExecutionRepo.FindAllByEncounter(encounterID)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Executions not found"))
	}

	return executions, nil
}

func (service *EncounterExecutionService) GetAllBySocialEncounter(encounterID uint64) ([]model.EncounterExecution, error) {
	executions, err := service.ExecutionRepo.FindAllByType(encounterID, model.Social)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Executions not found"))
	}

	return executions, nil
}

func (service *EncounterExecutionService) GetAllByLocationEncounter(encounterID uint64) ([]model.EncounterExecution, error) {
	executions, err := service.ExecutionRepo.FindAllByType(encounterID, model.Location)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Executions not found"))
	}

	return executions, nil
}

func (service *EncounterExecutionService) GetByEncounterAndTourist(encounterID, touristID uint64) (*model.EncounterExecution, error) {
	execution, err := service.ExecutionRepo.FindByEncounterAndTourist(encounterID, touristID)
	if err != nil {
		return nil, fmt.Errorf("execution not found")
	}

	return execution, nil
}

func (service *EncounterExecutionService) UpdateRange(encounters []model.EncounterExecution) ([]model.EncounterExecution, error) {
	updatedExecutions, err := service.ExecutionRepo.UpdateRange(encounters)
	if err != nil {
		return nil, fmt.Errorf("failed to update executions: %v", err)
	}

	return updatedExecutions, nil
}

func (service *EncounterExecutionService) checkPermission(ID uint64, touristID uint64) error {
	execution, err := service.ExecutionRepo.FindByID(ID)
	if err != nil {
		return fmt.Errorf("failed to find execution: %v", err)
	}

	if execution.TouristID != touristID {
		return fmt.Errorf("tourist does not have permission")
	}

	return nil
}

// Activate encounter
func (service *EncounterExecutionService) Activate(encounterID, touristID uint64, touristLongitude, touristLatitude float64) (model.EncounterExecution, error) {
	execution, err := service.ExecutionRepo.FindByEncounterAndTourist(encounterID, touristID)
	if err != nil {
		return model.EncounterExecution{}, fmt.Errorf("execution not found")
	}

	isExecutionCompleted := execution.Status == model.Completed
	if isExecutionCompleted {
		return model.EncounterExecution{}, fmt.Errorf("execution is already completed")
	}

	isTouristInRange := service.isTouristInRange(*execution, touristLongitude, touristLatitude)

	if !isTouristInRange {
		return model.EncounterExecution{}, fmt.Errorf("tourist not in range")
	}

	execution.Activate()

	updatedExecution, err := service.ExecutionRepo.Update(*execution)
	if err != nil {
		return model.EncounterExecution{}, fmt.Errorf("execution cannot be updated: %v", err)
	}

	return updatedExecution, nil
}

func (service *EncounterExecutionService) Complete(executionID, touristID uint64, touristLongitude, touristLatitude float64) (*model.EncounterExecution, int32, error) {
	execution, err := service.ExecutionRepo.FindByID(executionID)
	if err != nil {
		return nil, 0, fmt.Errorf("execution with ID %d not found", executionID)
	}

	if execution.TouristID != touristID {
		return nil, 0, fmt.Errorf("not tourist encounter execution")
	}

	if execution.Status != model.Active {
		return nil, 0, fmt.Errorf("not valid status")
	}

	isInRange := service.isTouristInRange(*execution, touristLongitude, touristLatitude)

	if !isInRange {
		return nil, 0, fmt.Errorf("tourist not in range")
	}

	execution.Complete()

	if execution.Encounter.Type == model.Social {
		err := service.updateAllCompletedSocial(execution.EncounterID)
		if err != nil {
			return nil, 0, err
		}
	}

	if execution.Encounter.Type == model.Location {
		err := service.updateAllCompletedLocation(execution.EncounterID)
		if err != nil {
			return nil, 0, err
		}
	}

	updatedExecution, err := service.ExecutionRepo.Update(*execution)
	if err != nil {
		return &model.EncounterExecution{}, 0, fmt.Errorf("execution cannot be updated: %v", err)
	}

	// TODO update user XP points
	return &updatedExecution, updatedExecution.Encounter.XP, nil
}

func (service *EncounterExecutionService) GetVisibleByTour(touristID uint64, touristLongitude, touristLatitude float64, encounterIDs []uint64) (model.EncounterExecution, error) {
	// TODO get encounterIDs from checkpoints from front (InternalService)
	encounters, err := service.EncounterRepo.FindByIds(encounterIDs)
	if err != nil {
		return model.EncounterExecution{}, fmt.Errorf("encounters not found")
	}

	if len(encounters) == 0 {
		return model.EncounterExecution{}, fmt.Errorf("encounters not found")
	}

	var closestEncounter *model.Encounter

	for _, encounter := range encounters {
		if encounter.IsCloseEnough(touristLongitude, touristLatitude) {
			closestEncounter = &encounter
		}
	}

	if closestEncounter == nil {
		return model.EncounterExecution{}, fmt.Errorf("no near encounter")
	}

	bestDistance := closestEncounter.CalculateDistance(touristLongitude, touristLatitude)

	for _, encounter := range encounters {
		distance := encounter.CalculateDistance(touristLongitude, touristLatitude)
		isBetterDistance := distance < bestDistance
		isCloseEnough := encounter.IsCloseEnough(touristLongitude, touristLatitude)
		if isBetterDistance && isCloseEnough {
			bestDistance = distance
			closestEncounter = &encounter
		}

	}

	/*
			closestDistance := math.Inf(1) // Positive infinity

			for _, encounter := range encounters {
				if model.IsCloseEnough(encounter.Longitude, encounter.Latitude, touristLongitude, touristLatitude) {
					distance := model.CalculateDistance(encounter.Longitude, encounter.Latitude, touristLongitude, touristLatitude)
					if distance < closestDistance {
						closestEncounter = &encounter
						closestDistance = distance
					}
				}
			}


		if closestEncounter == nil {
			return model.EncounterExecution{}, errors.New("no close encounter found")
		}
	*/

	// TODO ask Anja what to do...
	//bestDistance := model.CalculateDistance(closestEncounter.Longitude, closestEncounter.Latitude, touristLongitude, touristLatitude)

	potentialEncounter, err := service.ExecutionRepo.FindByEncounterAndTourist(closestEncounter.ID, touristID)
	if err != nil {
		// If encounter execution not found, create a new one
		newEncounterExecution := model.EncounterExecution{
			EncounterID: closestEncounter.ID,
			Encounter:   *closestEncounter,
			TouristID:   touristID,
			Status:      model.Pending,
			StartTime:   time.Now(),
			EndTime:     time.Time{},
		}
		// Save the new encounter execution
		savedExecution, err := service.ExecutionRepo.Save(newEncounterExecution)

		if err != nil {
			return model.EncounterExecution{}, fmt.Errorf("execution cannot be created: %v", err)
		}
		return savedExecution, nil
	}

	return *potentialEncounter, nil

}

func (service *EncounterExecutionService) CheckIfInRange(executionID, touristID uint64, touristLongitude, touristLatitude float64) (*model.EncounterExecution, int32, error) {
	oldExecution, err := service.ExecutionRepo.FindByID(executionID)
	if err != nil {
		return nil, 0, fmt.Errorf("execution with ID %d not found", executionID)
	}
	if oldExecution.Status != model.Active {
		return nil, 0, fmt.Errorf("execution not activated")
	}

	socialEncounter, err := service.SocialEncounterRepo.FindById(oldExecution.EncounterID)
	if err != nil {
		return nil, 0, fmt.Errorf("encounter with ID %d not found", oldExecution.EncounterID)
	}

	socialEncounter.CheckIfInRange(touristLongitude, touristLatitude, touristID)

	_, err = service.SocialEncounterRepo.Update(*socialEncounter)

	if err != nil {
		return nil, 0, fmt.Errorf("execution cannot be updated: %v", err)
	}

	var newXP int32 = 0
	if socialEncounter.IsRequiredPeopleNumber() {
		socialExecutions, err := service.ExecutionRepo.FindAllByType(socialEncounter.EncounterID, model.Social)
		if err != nil {
			return nil, 0, fmt.Errorf(fmt.Sprintln("Executions not found"))
		}

		// TODO update XP for every finished encounter
		for _, activeSocial := range socialExecutions {
			if activeSocial.Status == model.Active && activeSocial.ID != executionID {
				_, XP, err := service.Complete(activeSocial.ID, activeSocial.TouristID, touristLatitude, touristLongitude)
				if err != nil {
					return nil, 0, fmt.Errorf("error completing execution: %v", err)
				}
				newXP += XP
			}
		}
		_, XP, err := service.Complete(executionID, touristID, touristLongitude, touristLatitude)
		if err != nil {
			return nil, 0, fmt.Errorf("error completing execution: %v", err)
		}
		newXP += XP
	}

	return oldExecution, newXP, nil

}

func (service *EncounterExecutionService) GetWithUpdatedLocation(encounterID, tourID, touristID uint64, touristLongitude, touristLatitude float64, encounterIDs []uint64) (*model.EncounterExecution, int32, error) {
	_, XP, err := service.CheckIfInRange(encounterID, touristID, touristLongitude, touristLatitude)
	if err != nil {
		return nil, 0, fmt.Errorf(fmt.Sprintln("not execution in range"))
	}

	encounter, err := service.GetVisibleByTour(tourID, touristLongitude, touristLatitude, encounterIDs)

	if err != nil {
		return nil, 0, fmt.Errorf(fmt.Sprintln("not visible execution for tour"))
	}

	return &encounter, XP, nil

}

func (service *EncounterExecutionService) GetHiddenLocationEncounterWithUpdatedLocation(encounterID, tourID, touristID uint64, touristLongitude, touristLatitude float64, encounterIDs []uint64) (*model.EncounterExecution, int32, error) {
	// TODO refactor function and extract logic
	_, XP, err := service.CheckIfInRangeLocation(encounterID, touristID, touristLongitude, touristLatitude)
	if err != nil {
		return nil, 0, fmt.Errorf(fmt.Sprintln("not execution in range"))
	}

	execution, err := service.GetVisibleByTour(tourID, touristLongitude, touristLatitude, encounterIDs)

	if err != nil {
		return nil, 0, fmt.Errorf(fmt.Sprintln("not visible execution for tour"))
	}

	return &execution, XP, nil

}

func (service *EncounterExecutionService) CheckIfInRangeLocation(executionID, touristID uint64, touristLongitude, touristLatitude float64) (*model.EncounterExecution, int32, error) {
	oldExecution, err := service.ExecutionRepo.FindByID(executionID)
	if err != nil {
		return nil, 0, fmt.Errorf("execution with ID %d not found", oldExecution.ID)
	}

	if oldExecution.Status != model.Active {
		return nil, 0, fmt.Errorf("execution is not active")
	}

	locationEncounter, err := service.LocationEncounterRepo.FindById(oldExecution.EncounterID)

	if err != nil {
		return nil, 0, fmt.Errorf("encounter with ID %d not found", locationEncounter.EncounterID)
	}

	isInRange := locationEncounter.CheckIfInRangeLocation(touristLongitude, touristLatitude)

	if !isInRange {
		return nil, 0, fmt.Errorf("encounter is not in range")
	}

	_, err = service.LocationEncounterRepo.Update(*locationEncounter)

	isLocationFound := locationEncounter.CheckIfLocationFound(touristLongitude, touristLatitude)
	if !isLocationFound {
		return nil, 0, fmt.Errorf("location is not found")
	}

	execution, XP, err := service.Complete(executionID, touristID, touristLongitude, touristLatitude)

	if err != nil {
		return oldExecution, 0, nil
	}

	return execution, XP, nil
}

func (service *EncounterExecutionService) updateAllCompletedSocial(encounterID uint64) error {
	executions, err := service.ExecutionRepo.FindAllByType(encounterID, model.Social)
	if err != nil {
		return err
	}

	_, err = service.ExecutionRepo.UpdateRange(executions)
	if err != nil {
		return err
	}

	return nil
}

func (service *EncounterExecutionService) updateAllCompletedLocation(encounterID uint64) error {
	executions, err := service.ExecutionRepo.FindAllByType(encounterID, model.Location)
	if err != nil {
		return err
	}

	_, err = service.ExecutionRepo.UpdateRange(executions)
	if err != nil {
		return err
	}

	return nil
}

func (service *EncounterExecutionService) isTouristInRange(execution model.EncounterExecution, touristLongitude, touristLatitude float64) bool {
	const thresholdDistance = 300

	switch execution.Encounter.Type {
	case model.Misc:
		distance := execution.Encounter.CalculateDistance(touristLongitude, touristLatitude)
		return distance < thresholdDistance

	case model.Social:
		socialEncounter, err := service.SocialEncounterRepo.FindById(execution.EncounterID)
		if err != nil {
			return false
		}
		socialEncDistance := execution.Encounter.CalculateDistance(touristLongitude, touristLatitude)
		return socialEncDistance <= socialEncounter.Range

	case model.Location:
		locationEncounter, err := service.LocationEncounterRepo.FindById(execution.EncounterID)
		if err != nil {
			return false
		}
		locationEncDistance := execution.Encounter.CalculateDistance(touristLongitude, touristLatitude)
		return locationEncDistance <= locationEncounter.Range

	default:
		return false
	}
}

func (service *EncounterExecutionService) GetActiveByTour(touristID uint64, encounterIDs []uint64) ([]model.EncounterExecution, error) {
	executions, err := service.ExecutionRepo.FindAllActiveByTourist(touristID)
	if err != nil {
		return nil, fmt.Errorf("executions not found")
	}

	filteredExecutions := make([]model.EncounterExecution, 0)
	encounterIDSet := make(map[uint64]bool)

	// Create a set of encounter IDs for efficient lookup
	for _, id := range encounterIDs {
		encounterIDSet[id] = true
	}

	// Filter executions based on encounter IDs
	for _, execution := range executions {
		if encounterIDSet[execution.EncounterID] {
			filteredExecutions = append(filteredExecutions, execution)
		}
	}

	return filteredExecutions, nil
}
