package model

import (
	"encoding/json"
	"errors"
	"time"
)

type ExecutionStatus int

const (
	Completed ExecutionStatus = iota
	Abandoned
	InProgress
)

var executionStatusStrings = [...]string{"Completed", "Abandoned", "InProgress"}

func (s ExecutionStatus) String() string {
	if s < Completed || s > InProgress {
		return "Unknown"
	}
	return executionStatusStrings[s]
}

func (s *ExecutionStatus) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}
	switch statusStr {
	case "Completed":
		*s = Completed
	case "Abandoned":
		*s = Abandoned
	case "InProgress":
		*s = InProgress
	default:
		return errors.New("invalid status")
	}
	return nil
}

func (s ExecutionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type TourExecution struct {
	ID                   uint64                 `json:"id" gorm:"primaryKey;autoIncrement"`
	TouristID            uint64                 `json:"touristId" gorm:"not null"`
	TourID               uint64                 `json:"tourId" gorm:"not null"`
	Start                time.Time              `json:"start"`
	LastActivity         time.Time              `json:"lastActivity"`
	ExecutionStatus      ExecutionStatus        `json:"executionStatus"`
	CompletedCheckpoints []CheckpointCompletion `json:"completedCheckpoints" gorm:"foreignKey:TourExecutionID"`
	Tour                 Tour                   `json:"tour" gorm:"foreignKey:TourID"`
}
