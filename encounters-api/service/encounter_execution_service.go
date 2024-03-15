package service

import (
	"encounters/model"
	"encounters/repo"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"time"
)

type EncounterExecutionService struct {
	ExecutionRepo *repo.EncounterExecutionRepository
	EncounterRepo *repo.EncounterRepository
}

func NewEncounterExecutionService(db *gorm.DB) *EncounterExecutionService {
	return &EncounterExecutionService{
		ExecutionRepo: repo.NewEncounterExecutionRepository(db),
		EncounterRepo: repo.NewEncounterRepositoryRepository(db),
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
	executions, err := service.ExecutionRepo.FindAllBySocialEncounter(encounterID)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintln("Executions not found"))
	}

	return executions, nil
}

func (service *EncounterExecutionService) GetAllByLocationEncounter(encounterID uint64) ([]model.EncounterExecution, error) {
	executions, err := service.ExecutionRepo.FindAllByLocationEncounter(encounterID)
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
	// TO-DO get encounterIDs from checkpoints from front (InternalService)
	encounters, err := service.EncounterRepo.FindByIds(encounterIDs)
	if err != nil {
		return model.EncounterExecution{}, fmt.Errorf("encounters not found")
	}

	var closestEncounter *model.Encounter
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

	// TODO ask Anja what to do...
	//bestDistance := model.CalculateDistance(closestEncounter.Longitude, closestEncounter.Latitude, touristLongitude, touristLatitude)

	potentialEncounter, err := service.ExecutionRepo.FindByEncounterAndTourist(closestEncounter.ID, touristID)
	if err != nil {
		// If encounter execution not found, create a new one
		newEncounterExecution := model.EncounterExecution{
			EncounterID: closestEncounter.ID,
			Encounter:   *closestEncounter,
			TouristID:   touristID,
			Status:      model.Pending, // Assuming Pending is the initial status
			StartTime:   time.Now(),    // Set the start time to the current time
			EndTime:     time.Time{},   // Zero time for EndTime as it's pending
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

func (service *EncounterExecutionService) CheckIfInRange(executionID, touristId uint64, touristLongitude, touristLatitude float64) (*model.EncounterExecution, error) {
	oldExecution, err := service.ExecutionRepo.FindByID(executionID)
	if err != nil {
		return nil, fmt.Errorf("execution with ID %d not found", executionID)
	}
	if oldExecution.Status != model.Active {
		return nil, fmt.Errorf("execution not activated")
	}

	// TODO use SocialEncounterRepository to get SocialEncounter
	//socialEncounter, err := service.EncounterRepo.FindByID(oldExecution.EncounterID)
	// TODO add CheckIfInRange in SocialEncounter model

	// TODO use SocialEncounterRepository to Update SocialEncounter

	// TODO Check RequiredPeopleNumber from SocialEncounter

	/*
		allActiveSocialEncounters, err := service.ExecutionRepo.FindAllBySocialEncounter(socialEncounter.ID)
		if err != nil {
			return nil, fmt.Errorf("execution with ID %d not found", executionID)
		}
	*/

	return nil, nil

}

func (service *EncounterExecutionService) updateAllCompletedSocial(encounterID uint64) error {
	executions, err := service.ExecutionRepo.FindAllBySocialEncounter(encounterID)
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
	executions, err := service.ExecutionRepo.FindAllByLocationEncounter(encounterID)
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
	distance := model.CalculateDistance(execution.Encounter.Longitude, execution.Encounter.Latitude, touristLongitude, touristLatitude)

	if execution.Encounter.Type == model.Misc && distance < thresholdDistance {
		return true
	}
	// TODO Check for social
	if execution.Encounter.Type == model.Social {
		return true
	}
	// TODO Check for location
	if execution.Encounter.Type == model.Location {
		return true
	}

	return false
}
