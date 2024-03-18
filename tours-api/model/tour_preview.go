package model

import (
	"github.com/lib/pq"
)

type TourPreview struct {
	ID                uint64            `json:"id"`
	AuthorID          uint64            `json:"authorId"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	DemandignessLevel DemandignessLevel `json:"demandignessLevel"`
	Price             float64           `json:"price"`
	Tags              pq.StringArray    `json:"tags"`
	Equipment         []Equipment       `json:"equipment"`
	Checkpoint        Checkpoint        `json:"checkpoint"`
	TourRatings       []TourRating      `json:"tourRating"`
}
