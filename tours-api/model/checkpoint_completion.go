package model

import (
	"time"
)

type CheckpointCompletion struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	TourExecutionID uint64    `json:"execution_id"`
	CheckpointID    uint64    `json:"checkpoint_id"`
	CompletionTime  time.Time `json:"completion_time"`
}
