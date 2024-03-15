package model

import (
	"time"
)

type CheckpointCompletion struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	TourExecutionID uint64    `json:"tour_execution_id" gorm:"not null"`
	CheckpointID    uint64    `json:"checkpoint_id" gorm:"not null"`
	CompletionTime  time.Time `json:"completion_time"`
}
