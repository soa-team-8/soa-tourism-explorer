package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type BigIntSlice []uint64

func (b *BigIntSlice) Scan(value interface{}) error {
	if value == nil {
		*b = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return b.parseBytes(v)
	case string:
		return b.parseBytes([]byte(v))
	default:
		return errors.New("unsupported Scan, storing driver.Value type " + fmt.Sprintf("%T", v) + " into type BigIntSlice")
	}
}

func (b *BigIntSlice) parseBytes(value []byte) error {
	// Trim curly braces and whitespace
	valueStr := strings.Trim(string(value), "{} ")
	if valueStr == "" {
		*b = nil
		return nil
	}

	// Split by commas
	strValues := strings.Split(valueStr, ",")

	for _, str := range strValues {
		// Trim whitespace
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}

		// Parse uint64
		intValue, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*b = append(*b, intValue)
	}
	return nil
}

func (b BigIntSlice) Value() (driver.Value, error) {
	if b == nil {
		return nil, nil
	}

	var strValues []string
	for _, val := range b {
		strValues = append(strValues, strconv.FormatUint(val, 10))
	}
	return "{" + strings.Join(strValues, ",") + "}", nil
}

type SocialEncounter struct {
	ID                uint64      `gorm:"primaryKey;autoIncrement"`
	Encounter         Encounter   `gorm:"-"`
	RequiredPeople    int         `json:"required_people"`
	Range             float64     `json:"range"`
	ActiveTouristsIds BigIntSlice `json:"active_tourists_ids,omitempty" gorm:"type:bigint[]"`
}

func (se *SocialEncounter) CheckIfInRange(touristLongitude, touristLatitude float64, touristId uint64) int {
	distance := math.Acos(math.Sin(math.Pi/180*(se.Encounter.Latitude))*math.Sin(math.Pi/180*touristLatitude)+math.Cos(math.Pi/180*se.Encounter.Latitude)*math.Cos(math.Pi/180*touristLatitude)*math.Cos(math.Pi/180*se.Encounter.Longitude-math.Pi/180*touristLongitude)) * 6371000
	if distance > se.Range {
		se.RemoveTourist(touristId)
		return 0
	} else {
		se.AddTourist(touristId)
		return len(se.ActiveTouristsIds)
	}
}

func (se *SocialEncounter) AddTourist(touristId uint64) {
	if !contains(se.ActiveTouristsIds, touristId) {
		se.ActiveTouristsIds = append(se.ActiveTouristsIds, touristId)
	}
}

func (se *SocialEncounter) RemoveTourist(touristId uint64) {
	index := indexOf(se.ActiveTouristsIds, touristId)
	if index != -1 {
		se.ActiveTouristsIds = append(se.ActiveTouristsIds[:index], se.ActiveTouristsIds[index+1:]...)
	}
}

func (se *SocialEncounter) IsRequiredPeopleNumber() bool {
	numberOfTourists := len(se.ActiveTouristsIds)
	if numberOfTourists >= se.RequiredPeople {
		se.ClearActiveTourists()
	}
	return numberOfTourists >= se.RequiredPeople
}

func (se *SocialEncounter) ClearActiveTourists() {
	se.ActiveTouristsIds = []uint64{}
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
