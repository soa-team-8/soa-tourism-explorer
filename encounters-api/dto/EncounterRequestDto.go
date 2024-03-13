package dto

import "encounters/model"

// EncounterRequestDto predstavlja DTO (Data Transfer Object) za EncounterRequest
type EncounterRequestDto struct {
	EncounterId uint64 `json:"encounterId"`
	TouristId   uint64 `json:"touristId"`
	Status      string `json:"status"`
}

func (e *EncounterRequestDto) ToReqModel() model.EncounterRequest {
	status := mapStringToStatusReq(e.Status)
	return model.EncounterRequest{
		EncounterId: e.EncounterId,
		TouristId:   e.TouristId,
		Status:      status,
	}
}

func ToDtoListReq(encounters []model.EncounterRequest) []EncounterRequestDto {
	encounterDtos := make([]EncounterRequestDto, len(encounters))
	for i, encounter := range encounters {
		encounterDtos[i] = ToDtoReq(encounter)
	}
	return encounterDtos
}

func ToDtoReq(encounter model.EncounterRequest) EncounterRequestDto {
	return EncounterRequestDto{
		EncounterId: encounter.EncounterId,
		TouristId:   encounter.TouristId,
		Status:      mapStatusToStringReq(encounter.Status),
	}
}

func mapStringToStatusReq(status string) model.RequestStatus {
	switch status {
	case "OnHold":
		return model.OnHold
	case "Accepted":
		return model.Accepted
	}
	return model.Rejected
}

func mapStatusToStringReq(status model.RequestStatus) string {
	switch status {
	case model.Rejected:
		return "Rejected"
	case model.Accepted:
		return "Accepted"
	case model.OnHold:
		return "OnHold"
	default:
		return "Unknown"
	}
}
