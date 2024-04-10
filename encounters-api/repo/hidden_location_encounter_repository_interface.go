package repo

import "encounters/model"

type HiddenLocationRepository interface {
	Save(hiddenLocationEncounter model.HiddenLocationEncounter) (model.HiddenLocationEncounter, error)
	FindById(id uint64) (*model.HiddenLocationEncounter, error)
	Update(hiddenLocationEncounter model.HiddenLocationEncounter) (model.HiddenLocationEncounter, error)
	DeleteById(id uint64) error
	FindAll() ([]model.HiddenLocationEncounter, error)
}
