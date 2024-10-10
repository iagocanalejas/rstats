package types

import (
	"github.com/iagocanalejas/rstats/internal/db"
)

type League struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Gender   *string `json:"gender,omitempty"`
	Category *string `json:"category,omitempty"`
}

func NewLeagueFromDB(from *db.LeagueRow) *League {
	var gender, category *string
	if from.Gender != nil && *from.Gender != "" {
		gender = from.Gender
	}
	if from.Category != nil && *from.Category != "" {
		category = from.Category
	}

	return &League{
		ID:       from.ID,
		Name:     from.Name,
		Symbol:   from.Symbol,
		Gender:   gender,
		Category: category,
	}
}
