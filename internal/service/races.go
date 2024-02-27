package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/iagocanalejas/regatas/internal/db"
	"github.com/iagocanalejas/regatas/internal/service/races"
)

func (s *Service) SearchRaces(keywords string) ([]races.Race, error) {
	// filters should be sent in <key>:<value>, ...
	filters, err := buildFilters(keywords)
	if err != nil {
		log.Println("error building filters: ", err)
		return nil, err
	}

	log.Printf("searching for races with: %+v", *filters)
	flatRaces, err := s.db.SearchRaces(filters)
	if err != nil {
		log.Println("error loading races: ", err)
		return nil, err
	}

	rs := make([]races.Race, len(flatRaces))
	for idx, race := range flatRaces {
		rs[idx] = *races.New(race)
	}
	return rs, nil
}

func buildFilters(k string) (*db.RaceFilters, error) {
	filter := db.RaceFilters{}
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
			filter.Year = year
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
