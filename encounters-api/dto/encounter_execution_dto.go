package dto

import "time"

type EncounterExecutionDto struct {
	ID          int64        `json:"id"`
	EncounterID int64        `json:"encounterId"`
	Encounter   EncounterDto `json:"encounterDto"`
	TouristID   int64        `json:"touristId"`
	Status      string       `json:"status"`
	StartTime   time.Time    `json:"startTime"`
}
