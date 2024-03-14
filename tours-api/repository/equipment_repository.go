package repository

import (
	"errors"
	"gorm.io/gorm"
	"tours/model"
)

type EquipmentRepository struct {
	DB *gorm.DB
}

func (repo *EquipmentRepository) Save(equipment model.Equipment) error {
	result := repo.DB.Create(&equipment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *EquipmentRepository) Delete(id uint64) error {
	result := repo.DB.Delete(&model.Equipment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *EquipmentRepository) Update(equipment model.Equipment) error {
	existingEquipment, err := repo.FindByID(equipment.ID)
	if err != nil {
		return err
	}
	if existingEquipment == nil {
		return errors.New("equipment not found")
	}

	if equipment.Name == existingEquipment.Name && equipment.Description == existingEquipment.Description {
		return nil
	}

	result := repo.DB.Save(&equipment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *EquipmentRepository) FindAll() ([]model.Equipment, error) {
	var equipment []model.Equipment
	if err := repo.DB.Find(&equipment).Error; err != nil {
		return nil, err
	}
	return equipment, nil
}

func (repo *EquipmentRepository) FindAllPaged(page, pageSize int) ([]model.Equipment, error) {
	var equipment []model.Equipment
	offset := (page - 1) * pageSize
	result := repo.DB.Offset(offset).Limit(pageSize).Find(&equipment)
	if result.Error != nil {
		return nil, result.Error
	}
	return equipment, nil
}

func (repo *EquipmentRepository) FindByID(id uint64) (*model.Equipment, error) {
	var equipment model.Equipment
	if err := repo.DB.First(&equipment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &equipment, nil
}

func (repo *EquipmentRepository) FindByTourID(tourID int) ([]model.Equipment, error) {
	var tour model.Tour
	if err := repo.DB.Preload("Equipment").First(&tour, tourID).Error; err != nil {
		return nil, err
	}
	return tour.Equipment, nil
}
