package model

import (
	"github.com/lib/pq"
	"time"
)

type TourRating struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Rating       uint64         `json:"rating" gorm:"not null"`
	Comment      string         `json:"comment"`
	TouristID    uint64         `json:"tourist_id" gorm:"not null"`
	TourID       uint64         `json:"tour_id" gorm:"not null"`
	TourDate     time.Time      `json:"tour_date"`
	CreationDate time.Time      `json:"creation_date"`
	ImageNames   pq.StringArray `json:"image_names" gorm:"type:text[]"`
}
