package repo

import (
	"encounters/model"
	"fmt"

	"gorm.io/gorm"
)

type EncounterExecutionRepository struct {
	DB *gorm.DB
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
	if err := r.DB.Find(&encounterExecutions).Error; err != nil {
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
