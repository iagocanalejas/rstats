package types

import "github.com/iagocanalejas/rstats/internal/db"

type Flag struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Edition *int16 `json:"edition,omitempty"`
}

func NewFlagFromDB(from *db.FlagRow, edition *int16) *Flag {
	return &Flag{
		ID:      from.ID,
		Name:    from.Name,
		Edition: edition,
	}
}
