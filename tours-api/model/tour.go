package model

import (
	"encoding/json"
	"errors"
	"github.com/lib/pq"
	"time"
)

type Status int

const (
	Draft Status = iota
	Published
	Archived
)

type DemandignessLevel int

const (
	Easy DemandignessLevel = iota
	Medium
	Hard
)

func (d *DemandignessLevel) UnmarshalJSON(data []byte) error {
	var levelStr string
	if err := json.Unmarshal(data, &levelStr); err != nil {
		return err
	}
	switch levelStr {
	case "Easy":
		*d = Easy
	case "Medium":
		*d = Medium
	case "Hard":
		*d = Hard
	default:
		return errors.New("invalid demandigness level")
	}
	return nil
}

type PublishedTour struct {
	PublishingDate time.Time `json:"publishing_date"`
}

type ArchivedTour struct {
	ArchivingDate time.Time `json:"archiving_date"`
}

type TransportationType int

const (
	Walking TransportationType = iota
	Driving
	Cycling
)

type TourTime struct {
	ID             uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	TimeInSeconds  uint64
	Distance       uint64
	Transportation TransportationType
}

type Tour struct {
	ID                uint64            `json:"id" gorm:"primaryKey;autoIncrement"`
	AuthorID          uint64            `json:"AuthorID" gorm:"not null"`
	Name              string            `json:"name" gorm:"unique;not null;check:name != ''"`
	Description       string            `json:"description" gorm:"not null;check:description != ''"`
	DemandignessLevel DemandignessLevel `json:"DemandignessLevel" gorm:"type:int"`
	Status            Status            `json:"tour_status"`
	Price             float64           `json:"price"`
	Tags              pq.StringArray    `json:"tags" gorm:"type:text[]"`
	PublishedTours    []PublishedTour   `json:"published_tours" gorm:"type:jsonb"`
	ArchivedTours     []ArchivedTour    `json:"archived_tours" gorm:"type:jsonb"`
	TourTimes         []TourTime        `json:"tour_times" gorm:"type:jsonb"`
	Closed            bool              `json:"closed"`
}
