package model

type Equipment struct {
	ID          uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"unique;not null;check:name != ''"`
	Description string `json:"description" gorm:"not null;check:description != ''"`
}
