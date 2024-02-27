package db

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

type Participant struct {
	ID int64 `db:"id"`

	Gender   string `db:"gender"`
	Category string `db:"category"`
	Distance int    `db:"distance"`

	ClubId      int64          `db:"club_id"`
	ClubName    string         `db:"club_name"`
	ClubRawName sql.NullString `db:"club_raw_name"`

	IsDisqualified bool `db:"disqualified"`

	Laps   pq.StringArray `db:"laps"`
	Lane   sql.NullInt64  `db:"lane"`
	Series sql.NullInt64  `db:"series"`
}

func (r *Repository) GetParticipantsByRaceID(raceID int64) ([]Participant, error) {
	query, args, err := sq.
		Select("p.id", "p.gender", "p.category", "p.distance", "p.laps", "p.lane", "p.series", "p.club_id as club_id", "e.name as club_name", "p.club_name as club_raw_name",
			"((SELECT count(*) FROM penalty pe WHERE pe.participant_id = p.id AND disqualification) > 0) as disqualified").
		From("participant p").
		LeftJoin("entity e ON p.club_id = e.id").
		Where(sq.Eq{"p.race_id": raceID}).
		OrderBy("p.laps[ARRAY_UPPER(p.laps, 1)] ASC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var participants []Participant
	err = r.db.Select(&participants, query, args...)
	if err != nil {
		return nil, err
	}

	return participants, nil
}
