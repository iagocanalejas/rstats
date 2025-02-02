package types

import (
	"github.com/iagocanalejas/rstats/internal/db"
)

type Participant struct {
	ID     int64 `json:"id"`
	RaceID int64 `json:"race_id"`

	Gender   string `json:"gender"`
	Category string `json:"category"`
	Distance int    `json:"distance"`

	Club *Entity `json:"club"`

	IsDisqualified bool `json:"disqualified"`

	Laps   *[]string `json:"laps"`
	Lane   *int16    `json:"lane"`
	Series *int16    `json:"series"`

	Speed *float64 `json:"speed"`
}

func NewParticipantFromDB(from *db.ParticipantRow) *Participant {
	club := NewEntityFromDB(&db.EntityRow{ID: from.ClubId, Name: from.ClubName}, (*[]string)(from.ClubRawNames))

	return &Participant{
		ID:     from.ID,
		RaceID: from.RaceID,

		Gender:   from.Gender,
		Category: from.Category,
		Distance: from.Distance,

		Club: club,

		IsDisqualified: false,

		Laps:   (*[]string)(from.Laps),
		Lane:   from.Lane,
		Series: from.Series,

		Speed: nil,
	}
}

func NewParticipantWithSpeedFromDB(from *db.ParticipantRowWithSpeed) *Participant {
	club := NewEntityFromDB(&db.EntityRow{ID: from.ClubId, Name: from.ClubName}, (*[]string)(from.ClubRawNames))

	return &Participant{
		ID:     from.ID,
		RaceID: from.RaceID,

		Gender:   from.Gender,
		Category: from.Category,
		Distance: from.Distance,

		Club: club,

		IsDisqualified: false,

		Laps:   (*[]string)(from.Laps),
		Lane:   from.Lane,
		Series: from.Series,

		Speed: &from.Speed,
	}
}
