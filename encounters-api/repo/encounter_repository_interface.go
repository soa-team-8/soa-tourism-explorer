package repo

import "encounters/model"

type EncounterRepository interface {
	Save(encounter model.Encounter) (model.Encounter, error)
	FindByID(id uint64) (*model.Encounter, error)
	FindAll() ([]model.Encounter, error)
	DeleteByID(id uint64) error
	Update(encounter model.Encounter) (model.Encounter, error)
	FindByIds(ids []uint64) ([]model.Encounter, error)
}
