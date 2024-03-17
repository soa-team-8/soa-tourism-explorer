package model

import "math"

type EncounterStatus int

const (
	Draft EncounterStatus = iota
	Archived
	Published
)

type EncounterType int

const (
	Social EncounterType = iota
	Location
	Misc
)

type Encounter struct {
	ID          uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	AuthorID    uint64          `json:"authorId"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	XP          int32           `json:"XP"`
	Status      EncounterStatus `json:"status"`
	Type        EncounterType   `json:"type"`
	Longitude   float64         `json:"longitude"`
	Latitude    float64         `json:"latitude"`
	// List of changes
}

func (e *Encounter) MakeEncounterPublished() {
	e.Status = Published
}

func (e *Encounter) CalculateDistance(touristLongitude, touristLatitude float64) float64 {
	if e.Longitude == touristLongitude && e.Latitude == touristLatitude {
		return 0
	}

	eLonRad := e.degToRad(e.Longitude)
	eLatRad := e.degToRad(e.Latitude)
	touristLonRad := e.degToRad(touristLongitude)
	touristLatRad := e.degToRad(touristLatitude)

	// Calculate distance using Haversine formula
	deltaLon := touristLonRad - eLonRad
	deltaLat := touristLatRad - eLatRad
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(eLatRad)*math.Cos(touristLatRad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := c * 6371000 // Earth's radius in meters

	return distance
}

func (e *Encounter) degToRad(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

func (e *Encounter) IsCloseEnough(touristLongitude, touristLatitude float64) bool {
	const thresholdDistance = 1000 // 1000 meters or 1 kilometer

	distance := e.CalculateDistance(touristLongitude, touristLatitude)

	return distance <= thresholdDistance
}
