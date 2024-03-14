package model

import (
	"time"
)

type CheckpointCompletion struct {
	TourExecutionID int64
	CheckpointID    int64
	CompletionTime  time.Time
}
