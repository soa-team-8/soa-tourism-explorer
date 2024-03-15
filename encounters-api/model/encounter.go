package model

type EncounterStatus int

const (
	Draft EncounterStatus = iota
	Archived
	Published
)

type EncounterType int

const (
	Social EncounterType = iota
	Location
	Misc
)

type Encounter struct {
	ID          uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	AuthorID    uint64          `json:"authorId"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	XP          int32           `json:"XP"`
	Status      EncounterStatus `json:"status"`
	Type        EncounterType   `json:"type"`
	Longitude   float64         `json:"longitude"`
	Latitude    float64         `json:"latitude"`
	// List of changes
}

func (e *Encounter) MakeEncounterPublished() {
	e.Status = Published
}
