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

func (repo *TourExecutionRepository) FindInProgressByIds(uid int, tid int) (*model.TourExecution, error) {
	var tourExecution model.TourExecution
	if err := repo.DB.Where("tour_id = ? AND tourist_id = ? AND execution_status = ?", tid, uid, model.InProgress).First(&tourExecution).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tourExecution, nil
}

func (repo *TourExecutionRepository) FindByID(id int) (*model.TourExecution, error) {
	var tourExecution model.TourExecution
	if err := repo.DB.First(&tourExecution, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tourExecution, nil
}

func (repo *TourExecutionRepository) Update(tourExecution model.TourExecution) error {
	existingTourExecution, err := repo.FindByID(int(tourExecution.ID))
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
