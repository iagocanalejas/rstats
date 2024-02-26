package db

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/pgtype"
)

type Race struct {
	ID int64 `db:"id"`

	TrophyID      sql.NullInt64  `db:"trophy_id"`
	TrophyName    sql.NullString `db:"trophy_name"`
	TrophyEdition sql.NullInt64  `db:"trophy_edition"`

	FlagID      sql.NullInt64  `db:"flag_id"`
	FlagName    sql.NullString `db:"flag_name"`
	FlagEdition sql.NullInt64  `db:"flag_edition"`

	LeagueID     sql.NullInt64  `db:"league_id"`
	LeagueName   sql.NullString `db:"league_name"`
	LeagueGender sql.NullString `db:"league_gender"`

	Day  int64       `db:"day"`
	Date pgtype.Date `db:"date"`

	Gender   string `db:"gender"`
	Type     string `db:"type"`
	Modality string `db:"modality"`

	Laps        sql.NullInt64 `db:"laps"`
	Lanes       sql.NullInt64 `db:"lanes"`
	IsCancelled bool          `db:"cancelled"`

	Town    sql.NullString `db:"town"`
	Sponsor sql.NullString `db:"sponsor"`
}

func (r *Repository) SearchRaces(keywords string) ([]Race, error) {
	query, args, err := sq.
		Select("r.id", "r.day", "r.date", "r.gender", "r.type", "r.modality", "r.laps", "r.lanes", "r.cancelled", "r.town", "r.sponsor",
			"t.id as trophy_id", "t.name as trophy_name", "r.trophy_edition as trophy_edition",
			"f.id as flag_id", "f.name as flag_name", "r.flag_edition as flag_edition",
			"l.id as league_id", "l.name as league_name", "l.gender as league_gender",
		).
		From("race r").
		LeftJoin("trophy t ON t.id = r.trophy_id").
		LeftJoin("flag f ON f.id = r.flag_id").
		LeftJoin("league l ON l.id = r.league_id").
		Where(sq.Or{
			sq.ILike{"t.name": fmt.Sprint("%", keywords, "%")},
			sq.ILike{"f.name": fmt.Sprint("%", keywords, "%")},
			sq.ILike{"sponsor": fmt.Sprint("%", keywords, "%")},
		}).
		OrderBy("date DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var races []Race
	if err = r.db.Select(&races, query, args...); err != nil {
		return nil, err
	}

	return races, nil
}
