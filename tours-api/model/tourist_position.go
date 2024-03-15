package model

type TouristPosition struct {
	//ID        uint64  `json:"id"`
	CreatorID uint64  `json:"creatorId"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
