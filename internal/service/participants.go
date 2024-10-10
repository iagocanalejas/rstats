package service

import (
	"github.com/iagocanalejas/rstats/internal/db"
	"github.com/iagocanalejas/rstats/internal/types"
)

type GetYearSpeedsByParams struct {
	Club            *types.Entity
	League          *types.League
	Flag            *types.Flag
	Gender          string
	Category        string
	Day             int16
	Years           []int
	BranchTeams     bool
	OnlyLeagueRaces bool
	Normalize       bool
}

// GetYearSpeedsBy retrieves participant speeds grouped by year.
func (s *Service) GetYearSpeedsBy(params *GetYearSpeedsByParams) (*[]int, *map[int][]float64, error) {
	var clubID, leagueID, flagID int64
	if params.Club != nil {
		clubID = params.Club.ID
	}
	if params.League != nil {
		leagueID = params.League.ID
	}
	if params.Flag != nil {
		flagID = params.Flag.ID
	}

	return s.db.GetYearSpeedsBy(&db.GetYearSpeedsByParams{
		ClubID:          clubID,
		LeagueID:        leagueID,
		FlagID:          flagID,
		Gender:          params.Gender,
		Category:        params.Category,
		Day:             params.Day,
		Years:           params.Years,
		BranchTeams:     params.BranchTeams,
		OnlyLeagueRaces: params.OnlyLeagueRaces,
		Normalize:       params.Normalize,
	})
}

type GetNthSpeedsByParams struct {
	Index           int
	Club            *types.Entity
	League          *types.League
	Gender          string
	Category        string
	Day             int16
	Year            int16
	BranchTeams     bool
	OnlyLeagueRaces bool
	Normalize       bool
}

// GetNthSpeedsBy retrieves the nth fastest speeds for participants based on the provided filtering criteria.
func (s *Service) GetNthSpeedsBy(params *GetNthSpeedsByParams) (*[]float64, error) {
	var clubID, leagueID int64
	if params.Club != nil {
		clubID = params.Club.ID
	}
	if params.League != nil {
		leagueID = params.League.ID
	}

	return s.db.GetNthSpeedsBy(&db.GetNthSpeedsByParams{
		Index:           params.Index,
		ClubID:          clubID,
		LeagueID:        leagueID,
		Gender:          params.Gender,
		Category:        params.Category,
		Day:             params.Day,
		Year:            params.Year,
		BranchTeams:     params.BranchTeams,
		OnlyLeagueRaces: params.OnlyLeagueRaces,
		Normalize:       params.Normalize,
	})
}
