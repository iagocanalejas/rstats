package db

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
	"github.com/jackc/pgx/pgtype"
)

type RaceRow struct {
	ID int64 `db:"id"`

	TrophyID      *int64  `db:"trophy_id"`
	TrophyName    *string `db:"trophy_name"`
	TrophyEdition *int16  `db:"trophy_edition"`

	FlagID      *int64  `db:"flag_id"`
	FlagName    *string `db:"flag_name"`
	FlagEdition *int16  `db:"flag_edition"`

	LeagueID       *int64  `db:"league_id"`
	LeagueName     *string `db:"league_name"`
	LeagueGender   *string `db:"league_gender"`
	LeagueCategory *string `db:"league_category"`

	AssociatedID *int64 `db:"associated_id"`

	Day  int16       `db:"day"`
	Date pgtype.Date `db:"date"`

	Gender   string `db:"gender"`
	Type     string `db:"type"`
	Modality string `db:"modality"`

	Laps        *int16 `db:"laps"`
	Lanes       *int16 `db:"lanes"`
	Series      *int16 `db:"series"`
	IsCancelled bool   `db:"cancelled"`

	Sponsor *string `db:"sponsor"`

	Metadata []byte `db:"metadata"`
}

func (r *Repository) GetRaceByID(raceID int64) (*RaceRow, error) {
	query, args, err := sq.
		Select("r.id", "r.day", "r.date", "r.gender", "r.type", "r.modality", "r.laps", "r.lanes", "r.cancelled", "r.sponsor", "r.associated_id", "r.metadata",
			"t.id as trophy_id", "t.name as trophy_name", "r.trophy_edition",
			"f.id as flag_id", "f.name as flag_name", "r.flag_edition",
			"l.id as league_id", "l.name as league_name", "l.gender as league_gender", "l.category as league_category",
			"(SELECT MAX(series) FROM participant WHERE race_id = r.id) as series").
		From("race r").
		LeftJoin("trophy t ON r.trophy_id = t.id").
		LeftJoin("flag f ON r.flag_id = f.id").
		LeftJoin("league l ON r.league_id = l.id").
		Where(sq.Eq{"r.id": raceID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	assert.NoError(err, "building query", "query", query, "args", args)

	var race RaceRow
	if err = r.db.Get(&race, query, args...); err != nil {
		return nil, err
	}

	return &race, nil
}

type SearchRaceParams struct {
	Keywords      string
	Year          int16
	League        string
	LeagueID      int64
	Trophy        string
	TrophyID      int64
	Flag          string
	FlagID        int64
	Participant   string
	ParticipantID int64
}

func (r *Repository) SearchRaces(filters *SearchRaceParams) ([]RaceRow, error) {
	baseSelect := sq.
		Select("r.id", "r.day", "r.date", "r.gender", "r.type", "r.modality", "r.laps", "r.lanes", "r.cancelled", "r.sponsor",
			"t.id as trophy_id", "t.name as trophy_name", "r.trophy_edition as trophy_edition",
			"f.id as flag_id", "f.name as flag_name", "r.flag_edition as flag_edition",
			"l.id as league_id", "l.name as league_name", "l.gender as league_gender", "l.category as league_category",
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
	assert.NoError(err, "building query", "query", query, "args", args)

	var races []RaceRow
	if err = r.db.Select(&races, query, args...); err != nil {
		return nil, err
	}

	return races, nil
}
