package hn

import (
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Created       int64     `json:"created"`
	ParsedCreated time.Time `json:"-"`
	Karma         uint32    `json:"karma"`
	About         string    `json:"about,omitempty"`
	Submitted     []uint64  `json:"submitted,omitempty"`
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	var dup Alias

	if err := json.Unmarshal(data, &dup); err != nil {
		return fmt.Errorf("error unmarshaling user: %w", err)
	}

	*u = (User)(dup)
	u.ParsedCreated = time.Unix(u.Created, 0)

	return nil
}
