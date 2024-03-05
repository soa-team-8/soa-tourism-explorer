package encounter

import (
	"context"
	"encounters/model"

	"gorm.io/gorm"
)

type PostgresRepo struct {
	DB *gorm.DB
}

func (r *PostgresRepo) Save(ctx context.Context, encounter model.Encounter) error {
	result := r.DB.Create(&encounter)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgresRepo) FindByID(ctx context.Context, id uint64) (*model.Encounter, error) {
	var encounter model.Encounter

	if err := r.DB.First(&encounter, id).Error; err != nil {
		return nil, err
	}

	return &encounter, nil
}

func (r *PostgresRepo) DeleteByID(ctx context.Context, id uint64) error {
	result := r.DB.Delete(&model.Encounter{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PostgresRepo) FindAll(ctx context.Context) ([]model.Encounter, error) {
	var encounters []model.Encounter
	if err := r.DB.Find(&encounters).Error; err != nil {
		return nil, err
	}
	return encounters, nil
}

func (r *PostgresRepo) Update(ctx context.Context, encounter model.Encounter) error {
	result := r.DB.Save(&encounter)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
