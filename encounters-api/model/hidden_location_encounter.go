package model

import (
	"github.com/lib/pq"
)

type HiddenLocationEncounter struct {
	EncounterID       uint64 `gorm:"primaryKey;autoIncrement"`
	Encounter         Encounter
	LocationLongitude float64        `json:"location_longitude"`
	LocationLatitude  float64        `json:"location_latitude"`
	Image             pq.StringArray `json:"pictures" gorm:"type:text[]"`
}
