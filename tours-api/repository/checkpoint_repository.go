package repository

import (
	"errors"
	"gorm.io/gorm"
	"tours/model"
)

type CheckpointRepository struct {
	DB *gorm.DB
}

func (repo *CheckpointRepository) Save(checkpoint model.Checkpoint) (uint64, error) {
	result := repo.DB.Create(&checkpoint)
	if result.Error != nil {
		return 0, result.Error
	}
	return checkpoint.ID, nil
}

func (repo *CheckpointRepository) Delete(id uint64) error {
	result := repo.DB.Delete(&model.Checkpoint{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *CheckpointRepository) Update(checkpoint model.Checkpoint) error {
	existingCheckpoint, err := repo.FindByID(checkpoint.ID)
	if err != nil {
		return err
	}
	if existingCheckpoint == nil {
		return errors.New("checkpoint not found")
	}

	result := repo.DB.Save(&checkpoint)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *CheckpointRepository) FindByID(id uint64) (*model.Checkpoint, error) {
	var checkpoint model.Checkpoint
	if err := repo.DB.First(&checkpoint, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &checkpoint, nil
}

func (repo *CheckpointRepository) FindAll() ([]model.Checkpoint, error) {
	var checkpoint []model.Checkpoint
	if err := repo.DB.Find(&checkpoint).Error; err != nil {
		return nil, err
	}
	return checkpoint, nil
}

func (repo *CheckpointRepository) FindAllByTourID(tourID uint64) ([]model.Checkpoint, error) {
	var checkpoints []model.Checkpoint
	if err := repo.DB.Where("tour_id = ?", tourID).Find(&checkpoints).Error; err != nil {
		return nil, err
	}
	return checkpoints, nil
}

func (repo *CheckpointRepository) FindEncounterIDsByTour(tourID uint64) ([]uint64, error) {
	var encounterIDs []uint64

	// Retrieve checkpoints for the specified tourID
	checkpoints, err := repo.FindAllByTourID(tourID)
	if err != nil {
		return nil, err
	}

	// Extract encounter IDs from checkpoints
	for _, checkpoint := range checkpoints {
		if checkpoint.EncounterID != 0 {
			encounterIDs = append(encounterIDs, checkpoint.EncounterID)
		}
	}

	return encounterIDs, nil
}
