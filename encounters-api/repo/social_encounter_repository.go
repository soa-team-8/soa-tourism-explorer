package repo

import (
	"encounters/model"
	"gorm.io/gorm"
)

type SocialEncounterRepository struct {
	Db *gorm.DB
}

func NewSocialEncounterRepository(db *gorm.DB) *SocialEncounterRepository {
	return &SocialEncounterRepository{Db: db}
}

// CreateSocialEncounter creates a new social encounter record in the database
func (r *SocialEncounterRepository) Save(socialEncounter model.SocialEncounter) (model.SocialEncounter, error) {
	tx := r.Db.Begin()

	// Prvo upišite Encounter deo SocialEncountera
	if err := tx.Create(&socialEncounter.Encounter).Error; err != nil {
		tx.Rollback()
		return socialEncounter, err
	}

	// Zatim upišite SocialEncounter
	if err := tx.Create(&socialEncounter).Error; err != nil {
		tx.Rollback()
		return socialEncounter, err
	}

	// Commit transakcije ako nema grešaka
	if err := tx.Commit().Error; err != nil {
		return socialEncounter, err
	}

	return socialEncounter, nil
}

// GetSocialEncounterByID retrieves a social encounter record from the database by its ID
func (r *SocialEncounterRepository) FindById(id uint64) (*model.SocialEncounter, error) {
	socialEncounter := &model.SocialEncounter{}
	result := r.Db.First(socialEncounter, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return socialEncounter, nil
}

// UpdateSocialEncounter updates an existing social encounter record in the database
func (r *SocialEncounterRepository) Update(socialEncounter model.SocialEncounter) (model.SocialEncounter, error) {
	result := r.Db.Save(socialEncounter)
	if result.Error != nil {
		return model.SocialEncounter{}, result.Error
	}
	return socialEncounter, nil
}

// DeleteSocialEncounter deletes a social encounter record from the database by its ID
func (r *SocialEncounterRepository) DeleteById(id uint64) error {
	result := r.Db.Delete(&model.SocialEncounter{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindAll retrieves all social encounter records from the database
func (r *SocialEncounterRepository) FindAll() ([]model.SocialEncounter, error) {
	var socialEncounters []model.SocialEncounter
	result := r.Db.Find(&socialEncounters)
	if result.Error != nil {
		return nil, result.Error
	}
	return socialEncounters, nil
}
