package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/iagocanalejas/rstats/internal/db"
	"github.com/iagocanalejas/rstats/internal/types"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
	prettylog "github.com/iagocanalejas/rstats/internal/utils/pretty-log"
)

func (s *Service) GetRaceByID(raceID int64) (*types.Race, error) {
	dbRace, err := s.db.GetRaceByID(raceID)
	if err != nil {
		prettylog.Error("error loading race: %v", err)
		return nil, err
	}

	dbParticipants, err := s.db.GetParticipantsByRaceID(raceID)
	if err != nil {
		prettylog.Error("error loading participants: %v", err)
		return nil, err
	}

	// TODO: implement lap normalizations
	ps := make([]types.Participant, len(dbParticipants))
	for idx, participant := range dbParticipants {
		ps[idx] = *types.NewParticipantFromDB(&participant)
	}

	r := types.NewRaceFromDB(dbRace)
	r.Participants = ps

	return r, nil
}

func (s *Service) SearchRaces(keywords string) ([]types.Race, error) {
	// filters should be sent in <key>:<value>, ...
	filters, err := buildFilters(keywords)
	assert.NoError(err, "building filters keywords=%s", keywords)

	prettylog.Debug("searching races with filters=%v", *filters)
	flatRaces, err := s.db.SearchRaces(filters)
	if err != nil {
		prettylog.Error("error searching races: %v", err)
		return nil, err
	}

	rs := make([]types.Race, len(flatRaces))
	for idx, race := range flatRaces {
		rs[idx] = *types.NewRaceFromDB(&race)
	}
	return rs, nil
}

func buildFilters(k string) (*db.SearchRaceParams, error) {
	filter := db.SearchRaceParams{}
	for _, part := range strings.Split(k, ",") {
		if !strings.Contains(part, ":") {
			filter.Keywords = strings.TrimSpace(part)
			continue
		}

		keyValue := strings.Split(part, ":")
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid filter format: %s", part)
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		switch key {
		case "year":
			year, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing number from: %s", value)
			}
			filter.Year = int16(year)
		case "flag":
			filter.Flag = value
		case "flag_id":
			flagID, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing number from: %s", value)
			}
			filter.FlagID = flagID
		case "trophy":
			filter.Trophy = value
		case "trophy_id":
			trophyID, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing number from: %s", value)
			}
			filter.TrophyID = trophyID
		case "league":
			filter.League = value
		case "league_id":
			leagueID, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing number from: %s", value)
			}
			filter.LeagueID = leagueID
		case "participant":
			filter.Participant = value
		case "participant_id":
			participantID, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing number from: %s", value)
			}
			filter.ParticipantID = participantID
		default:
			return nil, fmt.Errorf("unknown filter key: %s", key)
		}
	}
	return &filter, nil
}
