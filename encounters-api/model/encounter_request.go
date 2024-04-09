package model

type RequestStatus int

const (
	OnHold RequestStatus = iota
	Accepted
	Rejected
)

type EncounterRequest struct {
	ID          uint64        `json:"id" gorm:"primaryKey;autoIncrement"`
	EncounterId uint64        `json:"encounterId"`
	TouristId   uint64        `json:"touristId"`
	Status      RequestStatus `json:"status"`
}

func (er *EncounterRequest) Accept() {
	er.Status = Accepted
}

func (er *EncounterRequest) Reject() {
	er.Status = Rejected
}
