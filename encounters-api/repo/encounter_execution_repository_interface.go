package repo

import "encounters/model"

type EncounterExecutionRepository interface {
	Save(encounterExecution model.EncounterExecution) (model.EncounterExecution, error)
	FindByID(id uint64) (*model.EncounterExecution, error)
	FindAll() ([]model.EncounterExecution, error)
	DeleteByID(id uint64) error
	Update(encounterExecution model.EncounterExecution) (model.EncounterExecution, error)

	FindAllByTourist(touristID uint64) ([]model.EncounterExecution, error)
	FindAllActiveByTourist(touristID uint64) ([]model.EncounterExecution, error)
	FindAllCompletedByTourist(touristID uint64) ([]model.EncounterExecution, error)

	FindByEncounter(encounterID uint64) (*model.EncounterExecution, error)
	FindAllByEncounter(encounterID uint64) ([]model.EncounterExecution, error)
	FindAllByType(encounterID uint64, encounterType model.EncounterType) ([]model.EncounterExecution, error)
	FindByEncounterAndTourist(encounterID, touristID uint64) (*model.EncounterExecution, error)

	UpdateRange(encounters []model.EncounterExecution) ([]model.EncounterExecution, error)
}
