package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/lib/pq"
)

type CheckpointSecret struct {
	Pictures    pq.StringArray `json:"pictures" gorm:"type:text[]"`
	Description string         `json:"description" gorm:"not null;check:description != ''"`
}

type Checkpoint struct {
	ID                    uint64           `json:"id" gorm:"primaryKey;autoIncrement"`
	TourID                uint64           `json:"tour_id" gorm:"not null"`
	AuthorID              uint64           `json:"author_id" gorm:"not null"`
	Longitude             float64          `json:"longitude" gorm:"not null"`
	Latitude              float64          `json:"latitude" gorm:"not null"`
	Name                  string           `json:"name" gorm:"not null;check:name != ''"`
	Description           string           `json:"description" gorm:"not null;check:description != ''"`
	Pictures              pq.StringArray   `json:"pictures" gorm:"type:text[]"`
	RequiredTimeInSeconds int              `json:"required_time_in_seconds" gorm:"not null"`
	Secret                CheckpointSecret `json:"secret" gorm:"type:jsonb"`
}

func (s *CheckpointSecret) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		return json.Unmarshal(bytes, s)
	}
	return errors.New("failed to unmarshal CheckpointSecret value")
}

func (s CheckpointSecret) Value() (driver.Value, error) {
	return json.Marshal(s)
}
