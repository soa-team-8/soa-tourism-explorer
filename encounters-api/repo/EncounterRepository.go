package repo

import (
	"encounters/model"
	"fmt"

	"gorm.io/gorm"
)

type EncounterRepository struct {
	DB *gorm.DB
}

func (r *EncounterRepository) Save(encounter model.Encounter) (model.Encounter, error) {
	result := r.DB.Create(&encounter)
	if result.Error != nil {
		return model.Encounter{}, result.Error
	}
	return encounter, nil
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

func (r *EncounterRepository) Update(encounter model.Encounter) (model.Encounter, error) {
	result := r.DB.Model(&model.Encounter{}).Where("id = ?", encounter.ID).Updates(&encounter)
	if result.Error != nil {
		return model.Encounter{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Encounter{}, fmt.Errorf("encounter with ID %d does not exist", encounter.ID)
	}

	return encounter, nil
}
