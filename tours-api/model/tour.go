package model

import (
	"github.com/lib/pq"
)

type TourStatus int

const (
	Draft TourStatus = iota
	Published
	Archived
)

type DifficultyLevel int

const (
	Easy DifficultyLevel = iota
	Medium
	Hard
)

type Tour struct {
	ID              uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	AuthorID        uint64          `json:"author_id" gorm:"not null"`
	Name            string          `json:"name" gorm:"unique;not null;check:name != ''"`
	Description     string          `json:"description" gorm:"not null;check:description != ''"`
	DifficultyLevel DifficultyLevel `json:"difficulty_level"`
	TourStatus      TourStatus      `json:"tour_status"`
	Price           float64         `json:"price"`
	Tags            pq.StringArray  `gorm:"type:text[]"`
}
