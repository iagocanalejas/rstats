package races

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/iagocanalejas/regatas/internal/db"
	"github.com/iagocanalejas/regatas/internal/service/flags"
	"github.com/iagocanalejas/regatas/internal/service/leagues"
	"github.com/iagocanalejas/regatas/internal/service/trophies"
	"github.com/iagocanalejas/regatas/internal/utils"
)

type Race struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	Trophy *trophies.Trophy `json:"trophy"`
	Flag   *flags.Flag      `json:"flag"`
	League *leagues.League  `json:"league"`

	Day  int64  `json:"day"`
	Date string `json:"date"`

	Gender   string `json:"gender"`
	Type     string `json:"type"`
	Modality string `json:"modality"`

	Laps        sql.NullInt64 `json:"laps"`
	Lanes       sql.NullInt64 `json:"lanes"`
	IsCancelled bool          `json:"cancelled"`

	Town    sql.NullString `json:"town"`
	Sponsor sql.NullString `json:"sponsor"`
}

func New(from db.Race) *Race {
	var trophy *trophies.Trophy
	if from.TrophyID.Valid {
		trophy = &trophies.Trophy{ID: from.TrophyID.Int64, Name: from.TrophyName.String, Edition: from.TrophyEdition.Int64}
	}

	var flag *flags.Flag
	if from.FlagID.Valid {
		flag = &flags.Flag{ID: from.FlagID.Int64, Name: from.FlagName.String, Edition: from.FlagEdition.Int64}
	}

	var league *leagues.League
	if from.LeagueID.Valid {
		league = &leagues.League{ID: from.LeagueID.Int64, Name: from.LeagueName.String}
	}

	return &Race{
		ID:   from.ID,
		Name: buildRaceName(from, false),

		Trophy: trophy,
		Flag:   flag,
		League: league,

		Day:  from.Day,
		Date: from.Date.Time.Format("02-01-2006"),

		Gender:   from.Gender,
		Type:     from.Type,
		Modality: from.Modality,

		Laps:        from.Laps,
		Lanes:       from.Lanes,
		IsCancelled: from.IsCancelled,

		Town:    from.Town,
		Sponsor: from.Sponsor,
	}
}

func buildRaceName(race db.Race, associated bool) string {
	day := ""
	if (!associated && race.Day > 1) || (associated && race.Day == 1) {
		day = fmt.Sprintf("XORNADA %d", race.Day)
	}

	gender := ""
	if race.Gender == "FEMALE" || (race.LeagueGender.Valid && race.LeagueGender.String == "FEMALE") {
		gender = "(FEMENINA)"
	}

	trophy := ""
	if race.TrophyID.Valid && race.TrophyEdition.Int64 > 0 {
		trophy = fmt.Sprintf("%s - %s", utils.Int2Roman(race.TrophyEdition.Int64), race.TrophyName.String)
		trophy = strings.Replace(trophy, "(CLASIFICATORIA)", "", -1)
	}

	flag := ""
	if race.FlagID.Valid && race.FlagEdition.Int64 > 0 {
		flag = fmt.Sprintf("%s - %s", utils.Int2Roman(race.FlagEdition.Int64), race.FlagName.String)
	}

	sponsor := ""
	if race.Sponsor.Valid {
		sponsor = race.Sponsor.String
	}

	nameParts := []string{trophy, flag, sponsor}
	var filteredParts []string
	for _, part := range nameParts {
		if part != "" {
			filteredParts = append(filteredParts, part)
		}
	}
	name := strings.Join(filteredParts, " - ")

	return strings.TrimSpace(fmt.Sprintf("%s %s %s", name, day, gender))
}
