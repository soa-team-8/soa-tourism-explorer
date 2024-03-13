package model

type SocialEncounter struct {
	ID                uint64   `json:"id" gorm:"primaryKey;autoIncrement"`
	RequiredPeople    int      `json:"required_people"`
	Range             float64  `json:"range"`
	ActiveTouristsIds []uint64 `json:"active_tourists_ids,omitempty" gorm:"type:bigint[]"`
}
