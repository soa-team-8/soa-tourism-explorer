package service

import (
	"errors"
	"fmt"
	"tours/model"
	"tours/repository"
)

type CheckpointService struct {
	CheckpointRepository *repository.CheckpointRepository
}

func (service *CheckpointService) Create(checkpoint model.Checkpoint) (uint64, error) {
	id, err := service.CheckpointRepository.Save(checkpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to create checkpoint: %w", err)
	}
	return id, nil
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

func (service *CheckpointService) CreateOrUpdateCheckpointSecret(id uint64, secret model.CheckpointSecret) error {
	checkpoint, err := service.CheckpointRepository.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve checkpoint with ID %d: %w", id, err)
	}
	if checkpoint == nil {
		return errors.New("checkpoint not found")
	}

	checkpoint.CheckpointSecret = secret

	if err := service.CheckpointRepository.Update(*checkpoint); err != nil {
		return fmt.Errorf("failed to update checkpoint with ID %d: %w", id, err)
	}

	return nil
}

func (service *CheckpointService) SetCheckpointEncounter(id uint64, encId uint64, isSecretPrerequisite bool) error {
	checkpoint, err := service.CheckpointRepository.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve checkpoint with ID %d: %w", id, err)
	}
	checkpoint.EncounterID = encId
	checkpoint.IsSecretPrerequisite = isSecretPrerequisite
	_ = service.Update(*checkpoint)
	return nil
}

func (service *CheckpointService) GetEncounterIDsByTour(tourID uint64) ([]uint64, error) {
	encounterIDs, err := service.CheckpointRepository.FindEncounterIDsByTour(tourID)
	if err != nil {
		return nil, fmt.Errorf("failed to get encounterIDs for Tour with id %d: %w", tourID, err)
	}

	return encounterIDs, nil
}
