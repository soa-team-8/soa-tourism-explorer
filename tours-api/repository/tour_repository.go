package repository

import (
	"errors"
	"gorm.io/gorm"
	"tours/model"
)

type TourRepository struct {
	DB *gorm.DB
}

func (repo *TourRepository) Save(tour model.Tour) (uint64, error) {
	result := repo.DB.Create(&tour)
	if result.Error != nil {
		return 0, result.Error
	}
	return tour.ID, nil
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
	var tours []model.Tour
	if err := repo.DB.Preload("Equipment").Preload("Checkpoints").Find(&tours).Error; err != nil {
		return nil, err
	}
	return tours, nil
}

func (repo *TourRepository) FindByID(id uint64) (*model.Tour, error) {
	var tour model.Tour
	if err := repo.DB.Preload("Equipment").Preload("Checkpoints").First(&tour, id).Error; err != nil {
		return nil, err
	}
	return &tour, nil
}

func (repo *TourRepository) AddEquipmentToTour(tourID uint64, equipmentID uint64) error {
	err := repo.DB.Create(&model.TourEquipment{
		TourID:      tourID,
		EquipmentID: equipmentID,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *TourRepository) RemoveEquipmentFromTour(tourID uint64, equipmentID uint64) error {
	result := repo.DB.Delete(&model.TourEquipment{}, "tour_id = ? AND equipment_id = ?", tourID, equipmentID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no matching entry found to remove")
	}
	return nil
}

func (repo *TourRepository) FindByAuthorID(authorID uint64) ([]model.Tour, error) {
	var tours []model.Tour

	if err := repo.DB.Where("author_id = ?", authorID).Preload("Equipment").Preload("Checkpoints").Find(&tours).Error; err != nil {
		return nil, err
	}
	return tours, nil
}
