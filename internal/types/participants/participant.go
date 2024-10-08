package participants

import (
	"database/sql"

	"github.com/iagocanalejas/rstats/internal/db"
	"github.com/iagocanalejas/rstats/internal/types/entities"
)

type Participant struct {
	ID int64 `json:"id"`

	Gender   string `json:"gender"`
	Category string `json:"category"`
	Distance int    `json:"distance"`

	Club *entities.Entity `json:"club"`

	IsDisqualified bool `json:"disqualified"`

	Laps   []string      `json:"laps"`
	Lane   sql.NullInt64 `json:"lane"`
	Series sql.NullInt64 `json:"series"`
}

func New(from db.Participant) *Participant {
	club := &entities.Entity{
		ID:      from.ClubId,
		Name:    from.ClubName,
		RawName: from.ClubRawName,
	}

	return &Participant{
		ID: from.ID,

		Gender:   from.Gender,
		Category: from.Category,
		Distance: from.Distance,

		Club: club,

		IsDisqualified: false,

		Laps:   from.Laps,
		Lane:   from.Lane,
		Series: from.Series,
	}
}
