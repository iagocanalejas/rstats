package service

import (
	"log"

	"github.com/iagocanalejas/rstats/internal/db"
)

func (s *Service) GetSpeedAVGByClubID(clubID int64, filters *db.ParticipantByClubFilters) ([]db.YearSpeeds, error) {
	if filters == nil {
		filters = &db.ParticipantByClubFilters{
			Gender:          "MALE",
			OnlyLeagueRaces: true,
			BranchTeams:     false,
		}
	}

	speeds, err := s.db.GetYearSpeedsByClubID(clubID, *filters)
	if err != nil {
		log.Println("error loading race: ", err)
		return nil, err
	}

	return speeds, nil
}
