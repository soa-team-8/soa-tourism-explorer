package repo

import (
	"encounters/model"
	"fmt"

	"gorm.io/gorm"
)

type EncounterExecutionRepository struct {
	DB *gorm.DB
}

func NewEncounterExecutionRepository(db *gorm.DB) *EncounterExecutionRepository {
	return &EncounterExecutionRepository{DB: db}
}

func (r *EncounterExecutionRepository) Save(encounterExecution model.EncounterExecution) (model.EncounterExecution, error) {
	result := r.DB.Create(&encounterExecution)
	if result.Error != nil {
		return model.EncounterExecution{}, result.Error
	}
	return encounterExecution, nil
}

func (r *EncounterExecutionRepository) FindByID(id uint64) (*model.EncounterExecution, error) {
	var encounterExecution model.EncounterExecution

	if err := r.DB.Preload("Encounter").First(&encounterExecution, id).Error; err != nil {
		return nil, err
	}

	return &encounterExecution, nil
}

func (r *EncounterExecutionRepository) FindAll() ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution
	if err := r.DB.Preload("Encounter").Find(&encounterExecutions).Error; err != nil {
		return nil, err
	}
	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) DeleteByID(id uint64) error {
	result := r.DB.Delete(&model.EncounterExecution{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *EncounterExecutionRepository) Update(encounterExecution model.EncounterExecution) (model.EncounterExecution, error) {
	result := r.DB.Model(&model.EncounterExecution{}).Where("id = ?", encounterExecution.ID).Updates(&encounterExecution)
	if result.Error != nil {
		return model.EncounterExecution{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.EncounterExecution{}, fmt.Errorf("encounter execution with ID %d does not exist", encounterExecution.ID)
	}

	return encounterExecution, nil
}

// New methods with complex queries
func (r *EncounterExecutionRepository) GetAllByTourist(touristID uint64) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	// Query to fetch EncounterExecutions with associated Encounter and matching TouristID
	if err := r.DB.Preload("Encounter").Where("tourist_id = ?", touristID).Find(&encounterExecutions).Error; err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) GetAllActiveByTourist(touristID int64) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	// Query to fetch active EncounterExecutions with associated Encounter and matching TouristID
	if err := r.DB.Preload("Encounter").
		Where("tourist_id = ? AND status = ?", touristID, model.Active).
		Find(&encounterExecutions).Error; err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) GetAllCompletedByTourist(touristID uint64) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	// Query to fetch completed EncounterExecutions with associated Encounter and matching TouristID
	if err := r.DB.Preload("Encounter").Where("tourist_id = ? AND status = ?", touristID, model.Completed).Find(&encounterExecutions).Error; err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) GetByEncounter(encounterID int64) (*model.EncounterExecution, error) {
	var encounterExecution model.EncounterExecution

	// Query to fetch EncounterExecution with associated Encounter matching the provided encounterID
	if err := r.DB.Preload("Encounter").Where("encounter_id = ?", encounterID).First(&encounterExecution).Error; err != nil {
		return nil, err
	}

	return &encounterExecution, nil
}
func (r *EncounterExecutionRepository) GetAllByEncounter(encounterID int64) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	// Query to fetch EncounterExecutions with associated Encounter ID
	if err := r.DB.Preload("Encounter").
		Where("encounter_id = ?", encounterID).
		Find(&encounterExecutions).Error; err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) GetAllBySocialEncounter(socialEncounterID int64) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	// Query to fetch EncounterExecutions with associated Encounter and matching social Encounter ID and type
	if err := r.DB.Preload("Encounter").
		Where("encounter_id = ? AND encounter.type = ?", socialEncounterID, model.Social).
		Find(&encounterExecutions).Error; err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) GetAllByLocationEncounter(locationEncounterID int64) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	// Query to fetch EncounterExecutions with associated Encounter and matching location Encounter ID and type
	if err := r.DB.Preload("Encounter").
		Where("encounter_id = ? AND encounter.type = ?", locationEncounterID, model.Location).
		Find(&encounterExecutions).Error; err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) GetByEncounterAndTourist(touristID, encounterID int64) (*model.EncounterExecution, error) {
	var encounterExecution model.EncounterExecution

	// Query to fetch EncounterExecution with associated Encounter and matching TouristID and EncounterID
	if err := r.DB.Preload("Encounter").
		Where("tourist_id = ? AND encounter_id = ?", touristID, encounterID).
		First(&encounterExecution).Error; err != nil {
		return nil, err
	}

	return &encounterExecution, nil
}

func (r *EncounterExecutionRepository) UpdateRange(encounters []model.EncounterExecution) ([]model.EncounterExecution, error) {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(&encounters).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return encounters, nil
}
