package types

import (
	"github.com/iagocanalejas/rstats/internal/db"
)

type Entity struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	RawName *[]string `json:"raw_name,omitempty"`
}

func NewEntityFromDB(from *db.EntityRow, rawNames *[]string) *Entity {
	return &Entity{
		ID:      from.ID,
		Name:    from.Name,
		RawName: rawNames,
	}
}
