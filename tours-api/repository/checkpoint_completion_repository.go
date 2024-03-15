package repository

import (
	"errors"
	"gorm.io/gorm"
	"tours/model"
)

type CheckpointCompletionRepository struct {
	DB *gorm.DB
}

func (repo *CheckpointCompletionRepository) Save(checkpointCompletion model.CheckpointCompletion) error {
	result := repo.DB.Create(&checkpointCompletion)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *CheckpointCompletionRepository) FindByIds(eid int, cid int) (*model.CheckpointCompletion, error) {
	var checkpointCompletion model.CheckpointCompletion
	if err := repo.DB.Where("tour_execution_id = ? AND checkpoint_id = ?", eid, cid).First(&checkpointCompletion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &checkpointCompletion, nil
}
