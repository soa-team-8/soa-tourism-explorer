package model

// RequestStatus je enumeracija za status zahteva
type RequestStatus int

const (
	OnHold RequestStatus = iota
	Accepted
	Rejected
)

// EncounterRequest je struktura koja predstavlja zahtev za susret
type EncounterRequest struct {
	ID          uint64        `json:"id" gorm:"primaryKey;autoIncrement"`
	EncounterId uint64        `json:"encounterId"`
	TouristId   uint64        `json:"touristId"`
	Status      RequestStatus `json:"status"`
}

// AcceptRequest postavlja status zahteva na Accepted
func (er *EncounterRequest) AcceptRequest() {
	er.Status = Accepted
}

// RejectRequest postavlja status zahteva na Rejected
func (er *EncounterRequest) RejectRequest() {
	er.Status = Rejected
}
