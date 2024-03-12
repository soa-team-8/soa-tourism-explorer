package repo

import (
	"encounters/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type EncounterRequestDatabaseRepository struct {
	db *gorm.DB
}

// NewEncounterRequestDatabaseRepository kreira novi EncounterRequestDatabaseRepository
func NewEncounterRequestDatabaseRepository(db *gorm.DB) *EncounterRequestDatabaseRepository {
	return &EncounterRequestDatabaseRepository{db: db}
}

// AcceptRequest prihvata zahtev za susret sa datim ID-om
func (r *EncounterRequestDatabaseRepository) AcceptRequest(id int) (*model.EncounterRequest, error) {
	requestToUpdate := &model.EncounterRequest{}
	err := r.db.First(requestToUpdate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Not found %d", id)
		}
		return nil, err
	}

	requestToUpdate.AcceptRequest()
	err = r.db.Save(requestToUpdate).Error
	if err != nil {
		return nil, err
	}

	return requestToUpdate, nil
}

// RejectRequest odbija zahtev za susret sa datim ID-om
func (r *EncounterRequestDatabaseRepository) RejectRequest(id int) (*model.EncounterRequest, error) {
	requestToUpdate := &model.EncounterRequest{}
	err := r.db.First(requestToUpdate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Not found %d", id)
		}
		return nil, err
	}

	requestToUpdate.RejectRequest()
	err = r.db.Save(requestToUpdate).Error
	if err != nil {
		return nil, err
	}

	return requestToUpdate, nil
}

func (r *EncounterRequestDatabaseRepository) Save(encounterReq model.EncounterRequest) (model.EncounterRequest, error) {
	result := r.db.Create(&encounterReq)
	if result.Error != nil {
		return model.EncounterRequest{}, result.Error
	}
	return encounterReq, nil
}
