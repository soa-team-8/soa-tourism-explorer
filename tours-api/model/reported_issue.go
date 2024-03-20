package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type ReportedIssueComment struct {
	Text         string    `json:"text"`
	CreationTime time.Time `json:"creationTime"`
	CreatorId    uint64    `json:"creatorId" gorm:"not null"`
}

func (r ReportedIssueComment) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *ReportedIssueComment) Scan(value interface{}) error {
	if value == nil {
		*r = ReportedIssueComment{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Scan source is not []byte")
	}
	return json.Unmarshal(bytes, r)
}

type ReportedIssue struct {
	ID          uint64                 `json:"id" gorm:"primaryKey;autoIncrement"`
	Category    string                 `json:"category"`
	Description string                 `json:"description"`
	Priority    uint64                 `json:"priority"`
	Time        time.Time              `json:"time"`
	TourID      uint64                 `json:"tourId" gorm:"not null"`
	Deadline    time.Time              `json:"deadline"`
	Closed      bool                   `json:"closed"`
	Resolved    bool                   `json:"resolved"`
	TouristID   uint64                 `json:"touristId" gorm:"not null"`
	Tour        Tour                   `json:"tour" gorm:"foreignKey:TourID"`
	Comments    []ReportedIssueComment `json:"comments" gorm:"type:jsonb"`
}
