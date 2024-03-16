package model

type TouristPosition struct {
	CreatorID uint64  `json:"creatorId"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
