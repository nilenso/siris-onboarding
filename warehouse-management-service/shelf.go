package wms

import (
	"errors"
)

type Shelf struct {
	Id           string `json:"id,omitempty"`
	Label        string `json:"label,omitempty"`
	Section      string `json:"section,omitempty"`
	Level        string `json:"level,omitempty"`
	ShelfBlockId string `json:"shelfBlockId,omitempty"`
}

var ShelfDoesNotExist = errors.New("shelf does not exist")
var InvalidShelfBlock = errors.New("invalid shelf block")

func NewShelf(label, section, level, shelfBlockId string) Shelf {
	return Shelf{
		Id:           generateUUID(),
		Label:        label,
		Section:      section,
		Level:        level,
		ShelfBlockId: shelfBlockId,
	}
}
