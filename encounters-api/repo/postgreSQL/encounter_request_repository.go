package postgreSQL

import (
	"encounters/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type EncounterRequestRepository struct {
	DB *gorm.DB
}

func NewEncounterRequestRepository(db *gorm.DB) *EncounterRequestRepository {
	return &EncounterRequestRepository{DB: db}
}

func (r *EncounterRequestRepository) AcceptRequest(id int) (*model.EncounterRequest, error) {
	requestToUpdate := &model.EncounterRequest{}
	err := r.DB.First(requestToUpdate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("not found %d", id)
		}
		return nil, err
	}

	requestToUpdate.Accept()
	err = r.DB.Save(requestToUpdate).Error
	if err != nil {
		return nil, err
	}

	return requestToUpdate, nil
}

func (r *EncounterRequestRepository) RejectRequest(id int) (*model.EncounterRequest, error) {
	requestToUpdate := &model.EncounterRequest{}
	err := r.DB.First(requestToUpdate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("not found %d", id)
		}
		return nil, err
	}

	requestToUpdate.Reject()
	err = r.DB.Save(requestToUpdate).Error
	if err != nil {
		return nil, err
	}

	return requestToUpdate, nil
}

func (r *EncounterRequestRepository) Save(encounterReq model.EncounterRequest) (model.EncounterRequest, error) {
	result := r.DB.Create(&encounterReq)
	if result.Error != nil {
		return model.EncounterRequest{}, result.Error
	}
	return encounterReq, nil
}

func (r *EncounterRequestRepository) FindAll() ([]model.EncounterRequest, error) {
	var encounterRequests []model.EncounterRequest
	if err := r.DB.Find(&encounterRequests).Error; err != nil {
		return nil, err
	}
	return encounterRequests, nil
}

func (r *EncounterRequestRepository) FindByID(id int) (*model.EncounterRequest, error) {
	var encounterRequest model.EncounterRequest
	if err := r.DB.First(&encounterRequest, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("not found %d", id)
		}
		return nil, err
	}
	return &encounterRequest, nil
}

func (r *EncounterRequestRepository) Update(encounterReq model.EncounterRequest) (*model.EncounterRequest, error) {
	err := r.DB.Save(&encounterReq).Error
	if err != nil {
		return nil, err
	}
	return &encounterReq, nil
}

func (r *EncounterRequestRepository) DeleteByID(id int) error {
	result := r.DB.Delete(&model.EncounterRequest{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
