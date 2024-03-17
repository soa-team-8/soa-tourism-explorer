package model

import (
	"math"
)

type SocialEncounter struct {
	EncounterID       uint64 `gorm:"primaryKey;autoIncrement"`
	Encounter         Encounter
	RequiredPeople    int       `json:"required_people"`
	Range             float64   `json:"range"`
	ActiveTouristsIds *[]uint64 `json:"active_tourists_ids,omitempty" gorm:"type:bigint[]"`
}

func (se *SocialEncounter) CheckIfInRange(touristLongitude, touristLatitude float64, touristId uint64) int {
	distance := math.Acos(math.Sin(math.Pi/180*(se.Encounter.Latitude))*math.Sin(math.Pi/180*touristLatitude)+math.Cos(math.Pi/180*se.Encounter.Latitude)*math.Cos(math.Pi/180*touristLatitude)*math.Cos(math.Pi/180*se.Encounter.Longitude-math.Pi/180*touristLongitude)) * 6371000
	if distance > se.Range {
		se.RemoveTourist(touristId)
		return 0
	} else {
		se.AddTourist(touristId)
		return len(*se.ActiveTouristsIds)
	}
}

func (se *SocialEncounter) AddTourist(touristId uint64) {
	if se.ActiveTouristsIds != nil && !contains(*se.ActiveTouristsIds, touristId) {
		*se.ActiveTouristsIds = append(*se.ActiveTouristsIds, touristId)
	}
}

func (se *SocialEncounter) RemoveTourist(touristId uint64) {
	if se.ActiveTouristsIds != nil && contains(*se.ActiveTouristsIds, touristId) {
		index := indexOf(*se.ActiveTouristsIds, uint64(touristId))
		*se.ActiveTouristsIds = append((*se.ActiveTouristsIds)[:index], (*se.ActiveTouristsIds)[index+1:]...)
	}
}

func (se *SocialEncounter) IsRequiredPeopleNumber() bool {
	numberOfTourists := len(*se.ActiveTouristsIds)
	if numberOfTourists >= se.RequiredPeople {
		se.ClearActiveTourists()
	}
	return numberOfTourists >= se.RequiredPeople
}

func (se *SocialEncounter) ClearActiveTourists() {
	*se.ActiveTouristsIds = []uint64{}
}

func contains(ids []uint64, id uint64) bool {
	for _, v := range ids {
		if v == uint64(id) {
			return true
		}
	}
	return false
}

func indexOf(ids []uint64, id uint64) int {
	for i, v := range ids {
		if v == id {
			return i
		}
	}
	return -1
}
