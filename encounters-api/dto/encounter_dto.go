package dto

import "encounters/model"

type EncounterDto struct {
	AuthorID          uint64    `json:"authorId"`
	ID                uint64    `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	XP                int       `json:"XP"`
	Status            string    `json:"status"`
	Type              string    `json:"type"`
	Longitude         float64   `json:"longitude"`
	Latitude          float64   `json:"latitude"`
	LocationLongitude *float64  `json:"location_longitude,omitempty"`
	LocationLatitude  *float64  `json:"location_latitude,omitempty"`
	Image             *string   `json:"image,omitempty"`
	Range             *float64  `json:"range,omitempty"`
	RequiredPeople    *int      `json:"required_people,omitempty"`
	ActiveTouristsIDs *[]uint64 `json:"active_tourists_ids,omitempty" gorm:"type:bigint[]"`
}

func (e *EncounterDto) ToModel() model.Encounter {
	status := mapStringToStatus(e.Status)
	encounterType := mapStringToType(e.Type)

	return model.Encounter{
		ID:          e.ID,
		AuthorID:    e.AuthorID,
		Name:        e.Name,
		Description: e.Description,
		XP:          uint64(e.XP),
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
		XP:                int(encounter.XP),
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

//-------------------------------------- social encounter mapping ---------------------------------------

func (e *EncounterDto) ToSocialModel() model.SocialEncounter {
	status := mapStringToStatus(e.Status)
	encounterType := mapStringToType(e.Type)

	return model.SocialEncounter{
		EncounterID: e.ID,
		Encounter: model.Encounter{
			ID:          e.ID,
			AuthorID:    e.AuthorID,
			Name:        e.Name,
			Description: e.Description,
			XP:          uint64(e.XP),
			Status:      status,
			Type:        encounterType,
			Longitude:   e.Longitude,
			Latitude:    e.Latitude,
		},
		RequiredPeople:    *e.RequiredPeople,
		Range:             *e.Range,
		ActiveTouristsIds: e.ActiveTouristsIDs,
	}
}

func ToSocialDtoList(socialEncounters []model.SocialEncounter) []EncounterDto {
	encounterDtos := make([]EncounterDto, len(socialEncounters))
	for i, encounter := range socialEncounters {
		encounterDtos[i] = ToDto(encounter.Encounter)
		encounterDtos[i].RequiredPeople = &encounter.RequiredPeople
		encounterDtos[i].Range = &encounter.Range
		encounterDtos[i].ActiveTouristsIDs = encounter.ActiveTouristsIds
	}
	return encounterDtos
}

func ToSocialDto(socialEncounter model.SocialEncounter) EncounterDto {
	dto := ToDto(socialEncounter.Encounter)
	dto.RequiredPeople = &socialEncounter.RequiredPeople
	dto.Range = &socialEncounter.Range
	dto.ActiveTouristsIDs = socialEncounter.ActiveTouristsIds
	return dto
}
