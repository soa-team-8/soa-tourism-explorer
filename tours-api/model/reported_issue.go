package model

import (
	"time"
)

type ReportedIssueComment struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ReportID     uint64    `json:"reportId"`
	Text         string    `json:"text"`
	CreationTime time.Time `json:"creationTime"`
	CreatorId    uint64    `json:"creatorId" gorm:"not null"`
}

type ReportedIssue struct {
	ID          uint64                 `json:"id" gorm:"primaryKey;autoIncrement"`
	Category    string                 `json:"category"`
	Description string                 `json:"description"`
	Priority    uint64                 `json:"priority"`
	Time        time.Time              `json:"time"`
	TourID      uint64                 `json:"tourId" gorm:"not null"`
	Deadline    *time.Time             `json:"deadline"`
	Closed      bool                   `json:"closed"`
	Resolved    bool                   `json:"resolved"`
	TouristID   uint64                 `json:"touristId" gorm:"not null"`
	Tour        Tour                   `json:"tour" gorm:"foreignKey:TourID"`
	Comments    []ReportedIssueComment `json:"comments" gorm:"foreignKey:ReportID"`
}
