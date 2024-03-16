package model

import (
	"github.com/lib/pq"
	"time"
)

type TourRating struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Rating       uint64         `json:"rating" gorm:"not null"`
	Comment      string         `json:"comment"`
	TouristID    uint64         `json:"touristId" gorm:"not null"`
	TourID       uint64         `json:"tourId" gorm:"not null"`
	TourDate     time.Time      `json:"tourDate"`
	CreationDate time.Time      `json:"creationDate"`
	ImageNames   pq.StringArray `json:"imageNames" gorm:"type:text[]"`
}
