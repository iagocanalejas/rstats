package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iagocanalejas/rstats/internal/db"
	"github.com/iagocanalejas/rstats/internal/utils"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
)

type RaceMetadata struct {
	Datasource []struct {
		DatasourceName *string           `json:"datasource_name"`
		RefId          *string           `json:"ref_id"`
		Values         map[string]string `json:"values"`
	} `json:"datasource"`
}

type Race struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	Trophy *Trophy `json:"trophy"`
	Flag   *Flag   `json:"flag"`
	League *League `json:"league"`

	Day  int16  `json:"day"`
	Date string `json:"date"`

	Gender   string `json:"gender"`
	Type     string `json:"type"`
	Modality string `json:"modality"`

	Laps        *int16 `json:"laps"`
	Lanes       *int16 `json:"lanes"`
	Series      *int16 `json:"series"`
	IsCancelled bool   `json:"cancelled"`

	Sponsor *string `json:"sponsor"`

	Metadata *RaceMetadata `json:"metadata"`

	Participants []Participant `json:"participants"`
}

func NewRaceFromDB(from *db.RaceRow) *Race {
	var trophy *Trophy
	if from.TrophyID != nil {
		assert.NotNil(from.TrophyName, "trophy name is nil")
		trophy = &Trophy{ID: *from.TrophyID, Name: *from.TrophyName, Edition: from.TrophyEdition}
	}

	var flag *Flag
	if from.FlagID != nil {
		assert.NotNil(from.FlagName, "flag name is nil")
		flag = &Flag{ID: *from.FlagID, Name: *from.FlagName, Edition: from.FlagEdition}
	}

	var league *League
	if from.LeagueID != nil {
		assert.NotNil(from.LeagueName, "league name is nil")
		league = &League{ID: *from.LeagueID, Name: *from.LeagueName, Gender: from.LeagueGender, Category: from.LeagueCategory}
	}

	var metadata *RaceMetadata
	if len(from.Metadata) > 0 {
		_ = json.Unmarshal([]byte(from.Metadata), &metadata)
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
		Series:      from.Series,
		IsCancelled: from.IsCancelled,

		Sponsor: from.Sponsor,

		Metadata: metadata,
	}
}

func buildRaceName(race *db.RaceRow, associated bool) string {
	day := ""
	if (!associated && race.Day > 1) || (associated && race.Day == 1) {
		day = fmt.Sprintf("XORNADA %d", race.Day)
	}

	gender := ""
	if race.Gender == "FEMALE" || (race.LeagueGender != nil && *race.LeagueGender == "FEMALE") {
		gender = "(FEMENINA)"
	}

	trophy := ""
	if race.TrophyID != nil && *race.TrophyEdition > 0 {
		trophy = fmt.Sprintf("%s - %s", utils.Int2Roman(*race.TrophyEdition), *race.TrophyName)
		trophy = strings.Replace(trophy, "(CLASIFICATORIA)", "", -1)
	}

	flag := ""
	if race.FlagID != nil && *race.FlagEdition > 0 {
		flag = fmt.Sprintf("%s - %s", utils.Int2Roman(*race.FlagEdition), *race.FlagName)
	}

	sponsor := ""
	if race.Sponsor != nil && *race.Sponsor != "" {
		sponsor = *race.Sponsor
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
