package dto

import "encounters/model"

type EncounterDto struct {
	AuthorID          int64    `json:"author_id"`
	ID                uint64   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	XP                int      `json:"XP"`
	Status            string   `json:"status"`
	Type              string   `json:"type"`
	Longitude         float64  `json:"longitude"`
	Latitude          float64  `json:"latitude"`
	LocationLongitude *float64 `json:"location_longitude,omitempty"`
	LocationLatitude  *float64 `json:"location_latitude,omitempty"`
	Image             *string  `json:"image,omitempty"`
	Range             *float64 `json:"range,omitempty"`
	RequiredPeople    *int     `json:"required_people,omitempty"`
	ActiveTouristsIDs []int    `json:"active_tourists_ids,omitempty"`
}

func (e *EncounterDto) ToModel() model.Encounter {
	status := mapStatus(e.Status)
	encounterType := mapType(e.Type)

	return model.Encounter{
		ID:          uint64(e.ID),
		AuthorID:    uint64(e.AuthorID),
		Name:        e.Name,
		Description: e.Description,
		XP:          uint64(e.XP),
		Status:      status,
		Type:        encounterType,
		Longitude:   e.Longitude,
		Latitude:    e.Latitude,
	}
}

func mapStatus(status string) model.EncounterStatus {
	switch status {
	case "Draft":
		return model.Draft
	case "Archived":
		return model.Archived
	case "Published":
		return model.Published
	}
	return model.Draft
}

func mapType(encounterType string) model.EncounterType {
	switch encounterType {
	case "Social":
		return model.Social
	case "Location":
		return model.Location
	case "Misc":
		return model.Misc
	}

	return model.Social
}
