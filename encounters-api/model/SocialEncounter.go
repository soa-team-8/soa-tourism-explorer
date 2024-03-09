package model

type SocialEncounter struct {
	EncounterID uint64
	Encounter   Encounter

	RequiredPeople    int     `json:"required_people"`
	Range             float64 `json:"range"`
	ActiveTouristsIds []int   `json:"active_tourists_ids,omitempty"`
}
