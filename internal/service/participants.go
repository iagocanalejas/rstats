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
	Day             int
	Years           []int
	BranchTeams     bool
	OnlyLeagueRaces bool
	Normalize       bool
}

// GetYearSpeedsBy retrieves participant speeds grouped by year.
//
// Params:
// - params (*GetYearSpeedsByParams): A struct containing the input parameters to filter and control the query:
//   - Club (*Club): A reference to the club entity (optional). If provided, the club ID is used for filtering.
//   - League (*League): A reference to the league entity (optional). If provided, the league ID is used for filtering.
//   - Flag (*Flag): A reference to a competition flag (optional). If provided, the flag ID is used for filtering.
//   - Gender (string): Filter by participant gender.
//   - Category (string): Filter by race category.
//   - Day (time.Time): The specific day to filter the races.
//   - Years ([]int): A list of years to filter races.
//   - BranchTeams (bool): Whether to include branch teams in the query.
//   - OnlyLeagueRaces (bool): Flag to filter only league races.
//   - Normalize (bool): Whether to normalize the speeds by filtering out outliers.
//
// Returns:
// - (*[]int): A pointer to a slice of integers representing the years for which data was retrieved.
// - (*map[int][]float64): A pointer to a map where the key is the year, and the value is a slice of speeds for that year.
// - (error): An error if there was an issue in fetching the data from the repository.
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
	Day             int
	Year            int
	BranchTeams     bool
	OnlyLeagueRaces bool
	Normalize       bool
}

// GetNthSpeedsBy retrieves the nth fastest speeds for participants based on the provided filtering criteria.
//
// Params:
// - params (*GetNthSpeedsByParams): A struct containing the input parameters to filter and control the query:
//   - Index (int): The position of the speed to retrieve (e.g., 1 for fastest, 2 for second fastest).
//   - Club (*Club): A reference to the club entity (optional). If provided, the club ID is used for filtering.
//   - League (*League): A reference to the league entity (optional). If provided, the league ID is used for filtering.
//   - Gender (string): Filter by participant gender.
//   - Category (string): Filter by race category.
//   - Day (time.Time): The specific day to filter the races.
//   - Year (int): The year to filter races.
//   - BranchTeams (bool): Whether to include branch teams in the query.
//   - OnlyLeagueRaces (bool): Flag to filter only league races.
//   - Normalize (bool): Whether to normalize speeds by filtering out outliers.
//
// Returns:
// - (*[]float64): A pointer to a slice of float64 representing the nth fastest speeds for each race that matches the filtering criteria.
// - (error): An error if there was an issue in fetching the data from the repository.
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
