package model

import "time"

type EncounterExecution struct {
	ID          uint64                   `json:"id" gorm:"primaryKey"`
	EncounterID int64                    `json:"encounterId"`
	Encounter   Encounter                `json:"encounter" gorm:"-"`
	TouristID   int64                    `json:"touristId"`
	Status      EncounterExecutionStatus `json:"status"`
	StartTime   time.Time                `json:"startTime" json:"-"`
	EndTime     time.Time                `json:"endTime" json:"-"`
}

func (ee *EncounterExecution) Activate() {
	ee.Status = Active
	ee.StartTime = time.Now()
}

func (ee *EncounterExecution) Abandon() {
	ee.Status = Abandoned
}

func (ee *EncounterExecution) Complete() {
	ee.Status = Completed
	ee.EndTime = time.Now()
}

type EncounterExecutionStatus int

const (
	Pending EncounterExecutionStatus = iota
	Completed
	Active
	Abandoned
)
