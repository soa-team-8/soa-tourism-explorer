package dto

import "encounters/model"

type EncounterDto struct {
	AuthorID          uint64   `json:"authorId"`
	ID                uint64   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	XP                int32    `json:"XP"`
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
	status := mapStringToStatus(e.Status)
	encounterType := mapStringToType(e.Type)

	return model.Encounter{
		ID:          e.ID,
		AuthorID:    e.AuthorID,
		Name:        e.Name,
		Description: e.Description,
		XP:          e.XP,
		Status:      status,
		Type:        encounterType,
		Longitude:   e.Longitude,
		Latitude:    e.Latitude,
	}
}

func ToDtoList(encounters []model.Encounter) []EncounterDto {
	encounterDtos := make([]EncounterDto, len(encounters))
	for i, encounter := range encounters {
		encounterDtos[i] = ToDto(encounter)
	}
	return encounterDtos
}

func ToDto(encounter model.Encounter) EncounterDto {
	return EncounterDto{
		AuthorID:          encounter.AuthorID,
		ID:                encounter.ID,
		Name:              encounter.Name,
		Description:       encounter.Description,
		XP:                encounter.XP,
		Status:            mapStatusToString(encounter.Status),
		Type:              mapTypeToString(encounter.Type),
		Longitude:         encounter.Longitude,
		Latitude:          encounter.Latitude,
		LocationLongitude: nil,
		LocationLatitude:  nil,
		Image:             nil,
		Range:             nil,
		RequiredPeople:    nil,
		ActiveTouristsIDs: nil,
	}
}

func mapStringToStatus(status string) model.EncounterStatus {
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

func mapStringToType(encounterType string) model.EncounterType {
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

func mapStatusToString(status model.EncounterStatus) string {
	switch status {
	case model.Draft:
		return "Draft"
	case model.Archived:
		return "Archived"
	case model.Published:
		return "Published"
	default:
		return "Unknown"
	}
}

func mapTypeToString(encounterType model.EncounterType) string {
	switch encounterType {
	case model.Social:
		return "Social"
	case model.Location:
		return "Location"
	case model.Misc:
		return "Misc"
	default:
		return "Unknown"
	}
}
