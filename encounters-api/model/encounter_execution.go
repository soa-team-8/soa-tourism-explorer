package model

import (
	"encoding/json"
	"errors"
	"math"
	"time"
)

type EncounterExecution struct {
	ID          uint64                   `json:"id" gorm:"primaryKey;autoIncrement"`
	EncounterID uint64                   `json:"encounterId" gorm:"not null"`
	Encounter   Encounter                `json:"encounter" gorm:"foreignKey:EncounterID"`
	TouristID   uint64                   `json:"touristId" gorm:"not null"`
	Status      EncounterExecutionStatus `json:"status" gorm:"not null"`
	StartTime   time.Time                `json:"startTime" json:"-" gorm:"not null"`
	EndTime     time.Time                `json:"endTime" json:"-" gorm:"not null"`
}

func (ee *EncounterExecution) Activate() {
	ee.Status = Active
	ee.StartTime = time.Now()
}

func (ee *EncounterExecution) Abandon() {
	ee.Status = Abandoned
}

func (ee *EncounterExecution) Complete() {
	ee.Status = Completed
	ee.EndTime = time.Now()
}

func CalculateDistance(encounterLongitude, encounterLatitude, touristLongitude, touristLatitude float64) float64 {
	const earthRadius = 6371000

	lon1 := degreesToRadians(encounterLongitude)
	lat1 := degreesToRadians(encounterLatitude)
	lon2 := degreesToRadians(touristLongitude)
	lat2 := degreesToRadians(touristLatitude)

	dlon := lon2 - lon1
	dlat := lat2 - lat1
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return distance
}

func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func IsCloseEnough(encounterLongitude, encounterLatitude, touristLongitude, touristLatitude float64) bool {
	const thresholdDistance = 1000 // 1000 meters or 1 kilometer

	distance := CalculateDistance(encounterLongitude, encounterLatitude, touristLongitude, touristLatitude)
	return distance <= thresholdDistance
}

type EncounterExecutionStatus int

const (
	Pending EncounterExecutionStatus = iota
	Completed
	Active
	Abandoned
)

var encounterExecutionStatusStrings = [...]string{"Pending", "Completed", "Active", "Abandoned"}

func (s EncounterExecutionStatus) String() string {
	if s < Pending || s > Abandoned {
		return "Unknown"
	}
	return encounterExecutionStatusStrings[s]
}

func (s *EncounterExecutionStatus) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}
	switch statusStr {
	case "Pending":
		*s = Pending
	case "Completed":
		*s = Completed
	case "Active":
		*s = Active
	case "Abandoned":
		*s = Abandoned
	default:
		return errors.New("invalid status")
	}
	return nil
}

func (s EncounterExecutionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
