package postgreSQL

import (
	"encounters/model"
	"fmt"
	"gorm.io/gorm"
)

type SocialEncounterRepository struct {
	Db *gorm.DB
}

func NewSocialEncounterRepository(db *gorm.DB) *SocialEncounterRepository {
	return &SocialEncounterRepository{Db: db}
}

func (r *SocialEncounterRepository) Save(socialEncounter model.SocialEncounter) (model.SocialEncounter, error) {
	tx := r.Db.Begin()

	if err := tx.Create(&socialEncounter.Encounter).Error; err != nil {
		tx.Rollback()
		return socialEncounter, err
	}

	if err := tx.Create(&socialEncounter).Error; err != nil {
		tx.Rollback()
		return socialEncounter, err
	}

	if err := tx.Commit().Error; err != nil {
		return socialEncounter, err
	}

	return socialEncounter, nil
}

func (r *SocialEncounterRepository) FindByID(id uint64) (*model.SocialEncounter, error) {
	socialEncounter := &model.SocialEncounter{}

	// Preload the Encounter data
	result := r.Db.Preload("Encounter").First(socialEncounter, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return socialEncounter, nil
}

func (r *SocialEncounterRepository) Update(socialEncounter model.SocialEncounter) (model.SocialEncounter, error) {
	// Assuming r.Db is a valid GORM DB connection
	result := r.Db.Model(&model.SocialEncounter{}).Where("id = ?", socialEncounter.ID).Updates(&socialEncounter)

	if result.Error != nil {
		return model.SocialEncounter{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.SocialEncounter{}, fmt.Errorf("social encounter with ID %d does not exist", socialEncounter.ID)
	}

	return socialEncounter, nil
}

func (r *SocialEncounterRepository) DeleteByID(id uint64) error {
	// Start a transaction
	tx := r.Db.Begin()

	// Find the SocialEncounter record to get its Encounter ID
	var socialEncounter model.SocialEncounter
	if err := tx.First(&socialEncounter, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&model.Encounter{}, socialEncounter.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&model.SocialEncounter{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
func (r *SocialEncounterRepository) FindAll() ([]model.SocialEncounter, error) {
	var socialEncounters []model.SocialEncounter
	result := r.Db.Preload("Encounter").Find(&socialEncounters)
	if result.Error != nil {
		return nil, result.Error
	}
	return socialEncounters, nil
}
