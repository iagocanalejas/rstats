package types

import "github.com/iagocanalejas/rstats/internal/db"

type Trophy struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Edition *int16 `json:"edition,omitempty"`
}

func NewTrophyFromDB(from *db.TrophyRow, edition *int16) *Trophy {
	return &Trophy{
		ID:      from.ID,
		Name:    from.Name,
		Edition: edition,
	}
}
