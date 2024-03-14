package service

import (
	"encounters/model"
	"encounters/repo"
	"fmt"
	"gorm.io/gorm"
)

type EncounterExecutionService struct {
	ExecutionRepo *repo.EncounterExecutionRepository
}

func NewEncounterExecutionService(db *gorm.DB) *EncounterExecutionService {
	return &EncounterExecutionService{
		ExecutionRepo: repo.NewEncounterExecutionRepository(db),
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

func (service *EncounterExecutionService) checkPermission(id uint64, touristID uint64) error {
	execution, err := service.ExecutionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find execution: %v", err)
	}

	if execution.TouristID != touristID {
		return fmt.Errorf("tourist does not have permission")
	}

	return nil
}
