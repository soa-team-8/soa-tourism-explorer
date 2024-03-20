package repository

import (
	"errors"
	"gorm.io/gorm"
	"tours/model"
)

type TourExecutionRepository struct {
	DB *gorm.DB
}

func (repo *TourExecutionRepository) Save(tourExecution model.TourExecution) error {
	result := repo.DB.Create(&tourExecution)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *TourExecutionRepository) FindInProgressByIds(userId uint64, tourId uint64) (*model.TourExecution, error) {
	var tourExecution model.TourExecution
	if err := repo.DB.Preload("Tour.Equipment").Preload("Tour.Checkpoints").Preload("Tour").Preload("CompletedCheckpoints").Where("tour_id = ? AND tourist_id = ? AND execution_status = ?", tourId, userId, model.InProgress).First(&tourExecution).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tourExecution, nil
}

func (repo *TourExecutionRepository) ExistsByIDs(userId uint64, tourId uint64) (bool, error) {
	var count int64
	if err := repo.DB.Model(&model.TourExecution{}).
		Where("tour_id = ? AND tourist_id = ? AND execution_status = ?", tourId, userId, model.InProgress).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *TourExecutionRepository) FindByIds(userId uint64, tourId uint64) (*model.TourExecution, error) {
	var tourExecution model.TourExecution
	if err := repo.DB.Preload("Tour.Equipment").Preload("Tour.Checkpoints").Preload("Tour").Preload("CompletedCheckpoints").Where("tour_id = ? AND tourist_id = ? AND execution_status != ?", tourId, userId, model.Abandoned).First(&tourExecution).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tourExecution, nil
}

func (repo *TourExecutionRepository) FindByID(id uint64) (*model.TourExecution, error) {
	var tourExecution model.TourExecution
	if err := repo.DB.Preload("Tour.Equipment").Preload("Tour.Checkpoints").Preload("Tour").Preload("CompletedCheckpoints").First(&tourExecution, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tourExecution, nil
}

func (repo *TourExecutionRepository) Update(tourExecution model.TourExecution) error {
	existingTourExecution, err := repo.FindByID(tourExecution.ID)
	if err != nil {
		return err
	}
	if existingTourExecution == nil {
		return errors.New("tourExecution not found")
	}

	result := repo.DB.Save(&tourExecution)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
