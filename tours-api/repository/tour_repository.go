package repository

import (
	"errors"
	"gorm.io/gorm"
	"tours/model"
)

type TourRepository struct {
	DB *gorm.DB
}

func (repo *TourRepository) Save(tour model.Tour) error {
	result := repo.DB.Create(&tour)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *TourRepository) Delete(id uint64) error {
	result := repo.DB.Delete(&model.Tour{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *TourRepository) Update(tour model.Tour) error {
	existingTour, err := repo.FindByID(tour.ID)
	if err != nil {
		return err
	}
	if existingTour == nil {
		return errors.New("tour not found")
	}

	result := repo.DB.Save(&tour)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *TourRepository) FindAll() ([]model.Tour, error) {
	var tour []model.Tour
	if err := repo.DB.Find(&tour).Error; err != nil {
		return nil, err
	}
	return tour, nil
}

func (repo *TourRepository) FindByID(id uint64) (*model.Tour, error) {
	var tour model.Tour
	if err := repo.DB.First(&tour, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tour, nil
}
