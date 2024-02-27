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

type RaceFilters struct {
	Keywords      string
	Year          int64
	League        string
	LeagueID      int64
	Trophy        string
	TrophyID      int64
	Flag          string
	FlagID        int64
	Participant   string
	ParticipantID int64
}

func (r *Repository) SearchRaces(filters *RaceFilters) ([]Race, error) {
	baseSelect := sq.
		Select("r.id", "r.day", "r.date", "r.gender", "r.type", "r.modality", "r.laps", "r.lanes", "r.cancelled", "r.town", "r.sponsor",
			"t.id as trophy_id", "t.name as trophy_name", "r.trophy_edition as trophy_edition",
			"f.id as flag_id", "f.name as flag_name", "r.flag_edition as flag_edition",
			"l.id as league_id", "l.name as league_name", "l.gender as league_gender",
		).
		From("race r").
		LeftJoin("trophy t ON t.id = r.trophy_id").
		LeftJoin("flag f ON f.id = r.flag_id").
		LeftJoin("league l ON l.id = r.league_id")

	if filters.Keywords != "" {
		baseSelect = baseSelect.Where(sq.Or{
			sq.ILike{"t.name": fmt.Sprint("%", filters.Keywords, "%")},
			sq.ILike{"f.name": fmt.Sprint("%", filters.Keywords, "%")},
			sq.ILike{"sponsor": fmt.Sprint("%", filters.Keywords, "%")},
		})
	}

	if filters.Year > 0 {
		baseSelect = baseSelect.Where(sq.Eq{"EXTRACT(YEAR FROM r.date)": filters.Year})
	}

	if filters.Trophy != "" {
		baseSelect = baseSelect.Where(sq.ILike{"t.name": fmt.Sprint("%", filters.Trophy, "%")})
	}

	if filters.TrophyID > 0 {
		baseSelect = baseSelect.Where(sq.And{
			sq.NotEq{"r.trophy_id": nil},
			sq.Eq{"r.trophy_id": filters.TrophyID},
		})
	}

	if filters.Flag != "" {
		baseSelect = baseSelect.Where(sq.ILike{"f.name": fmt.Sprint("%", filters.Flag, "%")})
	}

	if filters.FlagID > 0 {
		baseSelect = baseSelect.Where(sq.And{
			sq.NotEq{"r.flag_id": nil},
			sq.Eq{"r.flag_id": filters.FlagID},
		})
	}

	if filters.League != "" {
		baseSelect = baseSelect.Where(sq.Or{
			sq.ILike{"t.name": fmt.Sprint("%", filters.League, "%")},
			sq.ILike{"t.symbol": fmt.Sprint("%", filters.League, "%")},
		})
	}

	if filters.LeagueID > 0 {
		baseSelect = baseSelect.Where(sq.And{
			sq.NotEq{"r.league_id": nil},
			sq.Eq{"r.league_id": filters.LeagueID},
		})
	}

	if filters.Participant != "" {
		subQ := fmt.Sprintf("EXISTS(SELECT 1 FROM participant p JOIN entity e on p.club_id = e.id WHERE p.race_id = r.id AND e.name ILIKE '%%%s%%')", filters.Participant)
		baseSelect = baseSelect.Where(subQ)
	}

	if filters.ParticipantID > 0 {
		subQ := fmt.Sprintf("EXISTS(SELECT 1 FROM participant p WHERE p.race_id = r.id AND p.club_id = %d)", filters.ParticipantID)
		baseSelect = baseSelect.Where(subQ)
	}

	query, args, err := baseSelect.
		OrderBy("date DESC, league_id").
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
