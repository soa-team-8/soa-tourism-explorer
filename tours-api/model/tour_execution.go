package model

import (
	"time"
)

type ExecutionStatus int

const (
	Completed ExecutionStatus = iota
	Abandoned
	InProgress
)

type TourExecution struct {
	ID                   uint64                 `json:"id" gorm:"primaryKey;autoIncrement"`
	TouristID            uint64                 `json:"tourist_id"`
	TourID               uint64                 `json:"tour_id" gorm:"not null"`
	Start                time.Time              `json:"start"`
	LastActivity         time.Time              `json:"last_activity"`
	ExecutionStatus      ExecutionStatus        `json:"execution_status"`
	CompletedCheckpoints []CheckpointCompletion `json:"completed_checkpoints" gorm:"foreignKey:TourExecutionID"`
	Tour                 Tour                   `json:"tour" gorm:"foreignKey:TourID"`
}
