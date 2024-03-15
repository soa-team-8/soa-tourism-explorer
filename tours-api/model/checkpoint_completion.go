package model

import (
	"time"
)

type CheckpointCompletion struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	TourExecutionID uint64    `json:"tourExecutionId" gorm:"not null"`
	CheckpointID    uint64    `json:"checkpointId" gorm:"not null"`
	CompletionTime  time.Time `json:"completitionTime"`
}
