package repo

import "encounters/model"

type SocialEncounterRepository interface {
	Save(socialEncounter model.SocialEncounter) (model.SocialEncounter, error)
	FindByID(id uint64) (*model.SocialEncounter, error)
	Update(socialEncounter model.SocialEncounter) (model.SocialEncounter, error)
	DeleteByID(id uint64) error
	FindAll() ([]model.SocialEncounter, error)
}
