package repo

import (
	"encounters/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type EncounterRequestRepository struct {
	Db *gorm.DB
}

// AcceptRequest prihvata zahtev za susret sa datim ID-om
func (r *EncounterRequestRepository) AcceptRequest(id int) (*model.EncounterRequest, error) {
	requestToUpdate := &model.EncounterRequest{}
	err := r.Db.First(requestToUpdate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Not found %d", id)
		}
		return nil, err
	}

	requestToUpdate.AcceptRequest()
	err = r.Db.Save(requestToUpdate).Error
	if err != nil {
		return nil, err
	}

	return requestToUpdate, nil
}

// RejectRequest odbija zahtev za susret sa datim ID-om
func (r *EncounterRequestRepository) RejectRequest(id int) (*model.EncounterRequest, error) {
	requestToUpdate := &model.EncounterRequest{}
	err := r.Db.First(requestToUpdate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Not found %d", id)
		}
		return nil, err
	}

	requestToUpdate.RejectRequest()
	err = r.Db.Save(requestToUpdate).Error
	if err != nil {
		return nil, err
	}

	return requestToUpdate, nil
}

func (r *EncounterRequestRepository) Save(encounterReq model.EncounterRequest) (model.EncounterRequest, error) {
	result := r.Db.Create(&encounterReq)
	if result.Error != nil {
		return model.EncounterRequest{}, result.Error
	}
	return encounterReq, nil
}

func (r *EncounterRequestRepository) FindAll() ([]model.EncounterRequest, error) {
	var encounterRequests []model.EncounterRequest
	if err := r.Db.Find(&encounterRequests).Error; err != nil {
		return nil, err
	}
	return encounterRequests, nil
}
