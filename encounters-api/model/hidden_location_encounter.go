package model

import (
	"github.com/lib/pq"
	"math"
)

type HiddenLocationEncounter struct {
	EncounterID       uint64 `gorm:"primaryKey;autoIncrement"`
	Encounter         Encounter
	LocationLongitude float64        `json:"location_longitude"`
	LocationLatitude  float64        `json:"location_latitude"`
	Image             pq.StringArray `json:"pictures" gorm:"type:text[]"`
	Range             float64        `json:"range"`
}

func (hle *HiddenLocationEncounter) CheckIfInRangeLocation(touristLongitude, touristLatitude float64) bool {
	distance := math.Acos(math.Sin(math.Pi/180*(hle.LocationLatitude))*math.Sin(math.Pi/180*touristLatitude)+math.Cos(math.Pi/180*hle.LocationLatitude)*math.Cos(math.Pi/180*touristLatitude)*math.Cos(math.Pi/180*hle.LocationLongitude-math.Pi/180*touristLongitude)) * 6371000
	return distance <= hle.Range
}

// Check if all conditions for completing a hidden location encounter are met
func (hle *HiddenLocationEncounter) CheckIfLocationFound(touristLongitude, touristLatitude float64) bool {
	distance := math.Acos(math.Sin(math.Pi/180*(hle.LocationLatitude))*math.Sin(math.Pi/180*touristLatitude)+math.Cos(math.Pi/180*hle.LocationLatitude)*math.Cos(math.Pi/180*touristLatitude)*math.Cos(math.Pi/180*hle.LocationLongitude-math.Pi/180*touristLongitude)) * 6371000
	return distance <= 5.5
}
