package model

import (
	"encoding/json"
	"io"
)

type User struct {
	ID       uint64 `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type Users []*User

func (o *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}
