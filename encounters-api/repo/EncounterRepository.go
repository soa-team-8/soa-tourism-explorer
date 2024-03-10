package repo

import (
	"encounters/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type EncounterRepository struct {
	DB *gorm.DB
}

func (r *EncounterRepository) Save(encounter model.Encounter) error {
	result := r.DB.Create(&encounter)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *EncounterRepository) FindByID(id uint64) (*model.Encounter, error) {
	var encounter model.Encounter

	if err := r.DB.First(&encounter, id).Error; err != nil {
		return nil, err
	}

	return &encounter, nil
}

func (r *EncounterRepository) FindAll() ([]model.Encounter, error) {
	var encounters []model.Encounter
	if err := r.DB.Find(&encounters).Error; err != nil {
		return nil, err
	}
	return encounters, nil
}

func (r *EncounterRepository) DeleteByID(id uint64) error {
	result := r.DB.Delete(&model.Encounter{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *EncounterRepository) Update(encounter model.Encounter) error {
	// Check if the encounter exists in the database
	exists, err := r.isExist(encounter.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("encounter with ID %d does not exist", encounter.ID)
	}

	// Update the encounter
	result := r.DB.Save(&encounter)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *EncounterRepository) isExist(id uint64) (bool, error) {
	existingEncounter := model.Encounter{}
	err := r.DB.First(&existingEncounter, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
