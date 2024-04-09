package repo

import "encounters/model"

type EncounterRequestRepository interface {
	Save(encounterReq model.EncounterRequest) (model.EncounterRequest, error)
	FindAll() ([]model.EncounterRequest, error)
	FindByID(id int) (*model.EncounterRequest, error)
	Update(encounterReq model.EncounterRequest) (*model.EncounterRequest, error)
	DeleteByID(id int) error
}
