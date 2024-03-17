package model

import (
	"encoding/json"
	"errors"
	"time"
)

type EncounterExecution struct {
	ID          uint64                   `json:"id" gorm:"primaryKey;autoIncrement"`
	EncounterID uint64                   `json:"encounterId" gorm:"not null"`
	Encounter   Encounter                `json:"encounter" gorm:"foreignKey:EncounterID"`
	TouristID   uint64                   `json:"touristId" gorm:"not null"`
	Status      EncounterExecutionStatus `json:"status" gorm:"not null"`
	StartTime   time.Time                `json:"startTime" json:"-" gorm:"not null"`
	EndTime     time.Time                `json:"endTime" json:"-" gorm:"not null"`
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
	Active
	Completed
	Abandoned
)

var encounterExecutionStatusStrings = [...]string{"Pending", "Completed", "Active", "Abandoned"}

func (s *EncounterExecutionStatus) String() string {
	if *s < Pending || *s > Abandoned {
		return "Unknown"
	}
	return encounterExecutionStatusStrings[*s]
}

func (s *EncounterExecutionStatus) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}
	switch statusStr {
	case "Pending":
		*s = Pending
	case "Completed":
		*s = Completed
	case "Active":
		*s = Active
	case "Abandoned":
		*s = Abandoned
	default:
		return errors.New("invalid status")
	}
	return nil
}

func (s *EncounterExecutionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
