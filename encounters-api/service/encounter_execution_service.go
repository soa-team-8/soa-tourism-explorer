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
		ExecutionRepo: &repo.EncounterExecutionRepository{
			DB: db,
		},
	}
}

func (service *EncounterExecutionService) Create(execution model.EncounterExecution) (model.EncounterExecution, error) {
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

func (service *EncounterExecutionService) DeleteByID(id uint64) error {
	err := service.ExecutionRepo.DeleteByID(id)
	if err != nil {
		return fmt.Errorf("execution cannot be deleted: %v", err)
	}
	return nil
}

func (service *EncounterExecutionService) Update(execution model.EncounterExecution) (model.EncounterExecution, error) {
	updatedExecution, err := service.ExecutionRepo.Update(execution)
	if err != nil {
		return model.EncounterExecution{}, fmt.Errorf("execution cannot be updated: %v", err)
	}

	return updatedExecution, nil
}
