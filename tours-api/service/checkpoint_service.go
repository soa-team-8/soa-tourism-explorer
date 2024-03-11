package service

import (
	"fmt"
	"tours/model"
	"tours/repository"
)

type CheckpointService struct {
	CheckpointRepository *repository.CheckpointRepository
}

func (service *CheckpointService) Create(checkpoint model.Checkpoint) error {
	if err := service.CheckpointRepository.Save(checkpoint); err != nil {
		return fmt.Errorf("failed to create checkpoint: %w", err)
	}
	return nil
}

func (service *CheckpointService) Delete(id uint64) error {
	if err := service.CheckpointRepository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete checkpoint with ID %d: %w", id, err)
	}
	return nil
}

func (service *CheckpointService) Update(checkpoint model.Checkpoint) error {
	if err := service.CheckpointRepository.Update(checkpoint); err != nil {
		return fmt.Errorf("failed to update checkpoint with ID %d: %w", checkpoint.ID, err)
	}
	return nil
}

func (service *CheckpointService) GetAll() ([]model.Checkpoint, error) {
	checkpoint, err := service.CheckpointRepository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all checkpoints: %w", err)
	}
	return checkpoint, nil
}

func (service *CheckpointService) GetByID(id uint64) (*model.Checkpoint, error) {
	checkpoint, err := service.CheckpointRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkpoint with ID %d: %w", id, err)
	}
	return checkpoint, nil
}

func (service *CheckpointService) GetAllByTourID(tourID uint64) ([]model.Checkpoint, error) {
	checkpoints, err := service.CheckpointRepository.FindAllByTourID(tourID)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkpoints for Tour with id %d: %w", tourID, err)
	}
	return checkpoints, nil
}
