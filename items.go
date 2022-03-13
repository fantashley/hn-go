package hn

import (
	"encoding/json"
	"fmt"
	"time"
)

type Item struct {
	ID          uint64    `json:"id"`
	Deleted     bool      `json:"deleted,omitempty"`
	Type        ItemType  `json:"type,omitempty"`
	By          string    `json:"by,omitempty"`
	Time        int64     `json:"time,omitempty"`
	ParsedTime  time.Time `json:"-"`
	Text        string    `json:"text,omitempty"`
	Dead        bool      `json:"dead,omitempty"`
	Parent      uint64    `json:"parent,omitempty"`
	Poll        uint64    `json:"poll,omitempty"`
	Kids        []uint64  `json:"kids,omitempty"`
	URL         string    `json:"url,omitempty"`
	Score       int       `json:"score,omitempty"`
	Title       string    `json:"title,omitempty"`
	Parts       []uint64  `json:"parts,omitempty"`
	Descendants int       `json:"descendants,omitempty"`
}

type ItemType string

const (
	ItemTypeJob     = "job"
	ItemTypeStory   = "story"
	ItemTypeComment = "comment"
	ItemTypePoll    = "poll"
	ItemTypePollOpt = "pollopt"
)

type Story Item
type Comment Item
type Ask Item
type Job Item
type Poll Item
type PollOpt Item

func (i *Item) UnmarshalJSON(data []byte) error {
	type Alias Item
	var dup Alias

	if err := json.Unmarshal(data, &dup); err != nil {
		return fmt.Errorf("error unmarshaling item: %w", err)
	}

	*i = (Item)(dup)
	i.ParsedTime = time.Unix(i.Time, 0)

	return nil
}
