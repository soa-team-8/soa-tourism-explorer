package repo

import (
	"encounters/model"
	"gorm.io/gorm"
)

type HiddenLocationRepository struct {
	Db *gorm.DB
}

// CreateSocialEncounter creates a new social encounter record in the database
func (r *HiddenLocationRepository) Save(hiddenLocationEncounter model.HiddenLocationEncounter) (model.HiddenLocationEncounter, error) {
	tx := r.Db.Begin()

	// Prvo upišite Encounter deo SocialEncountera
	if err := tx.Create(&hiddenLocationEncounter.Encounter).Error; err != nil {
		tx.Rollback()
		return hiddenLocationEncounter, err
	}

	// Zatim upišite SocialEncounter
	if err := tx.Create(&hiddenLocationEncounter).Error; err != nil {
		tx.Rollback()
		return hiddenLocationEncounter, err
	}

	// Commit transakcije ako nema grešaka
	if err := tx.Commit().Error; err != nil {
		return hiddenLocationEncounter, err
	}

	return hiddenLocationEncounter, nil
}

// GetSocialEncounterByID retrieves a social encounter record from the database by its ID
func (r *HiddenLocationRepository) FindById(id uint64) (*model.HiddenLocationEncounter, error) {
	hiddenLocationEncounter := &model.HiddenLocationEncounter{}
	result := r.Db.First(hiddenLocationEncounter, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return hiddenLocationEncounter, nil
}

// UpdateSocialEncounter updates an existing social encounter record in the database
func (r *HiddenLocationRepository) Update(hiddenLocationEncounter model.HiddenLocationEncounter) (model.HiddenLocationEncounter, error) {
	result := r.Db.Save(hiddenLocationEncounter)
	if result.Error != nil {
		return model.HiddenLocationEncounter{}, result.Error
	}
	return hiddenLocationEncounter, nil
}

// DeleteSocialEncounter deletes a social encounter record from the database by its ID
func (r *HiddenLocationRepository) DeleteById(id uint64) error {
	result := r.Db.Delete(&model.HiddenLocationEncounter{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindAll retrieves all social encounter records from the database
func (r *HiddenLocationRepository) FindAll() ([]model.HiddenLocationEncounter, error) {
	var hiddenLocationEncounters []model.HiddenLocationEncounter
	result := r.Db.Find(&hiddenLocationEncounters)
	if result.Error != nil {
		return nil, result.Error
	}
	return hiddenLocationEncounters, nil
}
