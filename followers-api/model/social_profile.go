package model

import (
	"encoding/json"
	"io"
)

type SocialProfile struct {
	UserID       uint64    `json:"userId,omitempty"`
	FollowersIds []*uint64 `json:"followersIds,omitempty"`
	FollowedIds  []*uint64 `json:"followedIds,omitempty"`
}

func (o *SocialProfile) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *SocialProfile) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}
