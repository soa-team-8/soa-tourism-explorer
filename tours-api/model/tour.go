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

var statusStrings = [...]string{"Draft", "Published", "Archived"}

func (s Status) String() string {
	if s < Draft || s > Archived {
		return "Unknown"
	}
	return statusStrings[s]
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}
	switch statusStr {
	case "Draft":
		*s = Draft
	case "Published":
		*s = Published
	case "Archived":
		*s = Archived
	default:
		return errors.New("invalid status")
	}
	return nil
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type DemandignessLevel int

const (
	Easy DemandignessLevel = iota
	Medium
	Hard
)

var demandignessLevelStrings = [...]string{"Easy", "Medium", "Hard"}

func (d DemandignessLevel) String() string {
	if d < Easy || d > Hard {
		return "Unknown"
	}
	return demandignessLevelStrings[d]
}

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

func (d DemandignessLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

type PublishedTour struct {
	PublishingDate time.Time `json:"publishingDate"`
}

type ArchivedTour struct {
	ArchivingDate time.Time `json:"archivingDate"`
}

type TransportationType int

const (
	Walking TransportationType = iota
	Driving
	Cycling
)

var transportationTypeStrings = [...]string{"Walking", "Driving", "Cycling"}

func (d TransportationType) String() string {
	if d < Walking || d > Cycling {
		return "Unknown"
	}
	return transportationTypeStrings[d]
}

func (d *TransportationType) UnmarshalJSON(data []byte) error {
	var levelStr string
	if err := json.Unmarshal(data, &levelStr); err != nil {
		return err
	}
	switch levelStr {
	case "Walking":
		*d = Walking
	case "Driving":
		*d = Driving
	case "Cycling":
		*d = Cycling
	default:
		return errors.New("invalid transportation type level")
	}
	return nil
}

func (d TransportationType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

type TourTime struct {
	ID             uint64             `json:"id" gorm:"primaryKey;autoIncrement"`
	TimeInSeconds  uint64             `json:"timeInSeconds"`
	Distance       uint64             `json:"distance"`
	Transportation TransportationType `json:"transportation"`
}

type Tour struct {
	ID                uint64            `json:"id" gorm:"primaryKey;autoIncrement"`
	AuthorID          uint64            `json:"authorId" gorm:"not null"`
	Name              string            `json:"name" gorm:"unique;not null;check:name != ''"`
	Description       string            `json:"description" gorm:"not null;check:description != ''"`
	DemandignessLevel DemandignessLevel `json:"demandignessLevel" gorm:"type:int"`
	Status            Status            `json:"status"`
	Price             float64           `json:"price"`
	Tags              pq.StringArray    `json:"tags" gorm:"type:text[]"`
	PublishedTours    []PublishedTour   `json:"publishedTours" gorm:"type:jsonb"`
	ArchivedTours     []ArchivedTour    `json:"archivedTours" gorm:"type:jsonb"`
	Closed            bool              `json:"closed"`
	Equipment         []Equipment       `json:"equipment" gorm:"many2many:tour_equipments;"`
	Checkpoints       []Checkpoint      `json:"checkpoints" gorm:"foreignKey:TourID"`
}
