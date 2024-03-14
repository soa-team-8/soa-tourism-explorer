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

func NewEncounterRepositoryRepository(db *gorm.DB) *EncounterRepository {
	return &EncounterRepository{DB: db}
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

func (r *EncounterRepository) MakeEncounterPublished(id uint64) (*model.Encounter, error) {
	encounterToUpdate := &model.Encounter{}
	err := r.DB.First(encounterToUpdate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Not found %d", id)
		}
		return nil, err
	}

	encounterToUpdate.MakeEncounterPublished()
	err = r.DB.Save(encounterToUpdate).Error
	if err != nil {
		return nil, err
	}

	return encounterToUpdate, nil
}

func (r *EncounterRepository) FindByIds(ids []uint64) ([]model.Encounter, error) {
	var encounters []model.Encounter
	if len(ids) == 0 {
		return encounters, nil
	}

	query := r.DB.Where("id IN ?", ids).Find(&encounters)
	if query.Error != nil {
		return nil, query.Error
	}

	return encounters, nil
}
