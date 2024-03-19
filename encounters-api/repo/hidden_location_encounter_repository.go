package repo

import (
	"encounters/model"
	"fmt"
	"gorm.io/gorm"
)

type HiddenLocationRepository struct {
	Db *gorm.DB
}

func NewHiddenLocationRepository(db *gorm.DB) *HiddenLocationRepository {
	return &HiddenLocationRepository{Db: db}
}

func (r *HiddenLocationRepository) Save(hiddenLocationEncounter model.HiddenLocationEncounter) (model.HiddenLocationEncounter, error) {
	tx := r.Db.Begin()

	if err := tx.Create(&hiddenLocationEncounter.Encounter).Error; err != nil {
		tx.Rollback()
		return hiddenLocationEncounter, err
	}

	if err := tx.Create(&hiddenLocationEncounter).Error; err != nil {
		tx.Rollback()
		return hiddenLocationEncounter, err
	}

	if err := tx.Commit().Error; err != nil {
		return hiddenLocationEncounter, err
	}

	return hiddenLocationEncounter, nil
}

func (r *HiddenLocationRepository) FindById(id uint64) (*model.HiddenLocationEncounter, error) {
	hiddenLocationEncounter := &model.HiddenLocationEncounter{}

	result := r.Db.Preload("Encounter").First(hiddenLocationEncounter, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return hiddenLocationEncounter, nil
}

func (r *HiddenLocationRepository) Update(hiddenLocationEncounter model.HiddenLocationEncounter) (model.HiddenLocationEncounter, error) {
	result := r.Db.Model(&model.HiddenLocationEncounter{}).Where("encounter_id = ?", hiddenLocationEncounter.EncounterID).Updates(&hiddenLocationEncounter)

	if result.Error != nil {
		return model.HiddenLocationEncounter{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.HiddenLocationEncounter{}, fmt.Errorf("location encounter with ID %d does not exist", hiddenLocationEncounter.EncounterID)
	}

	return hiddenLocationEncounter, nil
}

/*
func (r *HiddenLocationRepository) Update(hiddenLocationEncounter model.HiddenLocationEncounter) (model.HiddenLocationEncounter, error) {
	result := r.Db.Save(hiddenLocationEncounter)
	if result.Error != nil {
		return model.HiddenLocationEncounter{}, result.Error
	}
	return hiddenLocationEncounter, nil
}
*/

/*
func (r *HiddenLocationRepository) DeleteById(id uint64) error {
	result := r.Db.Delete(&model.HiddenLocationEncounter{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
*/

func (r *HiddenLocationRepository) DeleteById(id uint64) error {
	tx := r.Db.Begin()

	if err := tx.Delete(&model.Encounter{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&model.HiddenLocationEncounter{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *HiddenLocationRepository) FindAll() ([]model.HiddenLocationEncounter, error) {
	var hiddenLocationEncounters []model.HiddenLocationEncounter
	result := r.Db.Preload("Encounter").Find(&hiddenLocationEncounters)
	if result.Error != nil {
		return nil, result.Error
	}
	return hiddenLocationEncounters, nil
}
