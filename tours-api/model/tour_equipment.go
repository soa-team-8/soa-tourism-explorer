package model

type TourEquipment struct {
	TourID      uint64 `gorm:"primaryKey"`
	EquipmentID uint64 `gorm:"primaryKey"`
}
