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
	ID                   uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TouristID            int             `json:"tourist_id"`
	TourID               int             `json:"tour_id"`
	Start                time.Time       `json:"start"`
	LastActivity         time.Time       `json:"last_activity"`
	ExecutionStatus      ExecutionStatus `json:"execution_status"`
	CompletedCheckpoints []CheckpointCompletion
}
