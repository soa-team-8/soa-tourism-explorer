package model

import (
	"encoding/json"
	"io"
)

type SocialProfile struct {
	UserId     uint64  `json:"id,omitempty"`
	Followers  []*User `json:"followers"`
	Followed   []*User `json:"followed"`
	Followable []*User `json:"followable"`
}

func (o *SocialProfile) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *SocialProfile) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}
