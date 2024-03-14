package model

type SocialEncounter struct {
	EncounterID       uint64 `gorm:"primaryKey;autoIncrement"`
	Encounter         Encounter
	RequiredPeople    int       `json:"required_people"`
	Range             float64   `json:"range"`
	ActiveTouristsIds *[]uint64 `json:"active_tourists_ids,omitempty" gorm:"type:bigint[]"`
}
