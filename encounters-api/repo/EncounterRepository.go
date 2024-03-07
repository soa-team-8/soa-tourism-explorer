package repo

import (
	"encounters/model"

	"gorm.io/gorm"
)

type EncounterRepository struct {
	DB *gorm.DB
}

func (r *EncounterRepository) Save(encounter model.Encounter) error {
	newID, err := r.generateID()
	if err != nil {
		return err
	}

	encounter.ID = newID

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
	result := r.DB.Save(&encounter)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *EncounterRepository) generateID() (uint64, error) {
	var maxID uint64
	if err := r.DB.Model(&model.Encounter{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID).Error; err != nil {
		return 0, err
	}
	return maxID + 1, nil
}
