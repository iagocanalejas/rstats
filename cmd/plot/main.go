package main

import (
	"fmt"
	"strings"

	"github.com/iagocanalejas/rstats/internal/service"
	"github.com/iagocanalejas/rstats/internal/types"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
	prettylog "github.com/iagocanalejas/rstats/internal/utils/pretty-log"
	"github.com/iagocanalejas/rstats/pkg/plotter"
	"github.com/spf13/pflag"
)

func main() {
	pflag.StringVarP(&plotType, "type", "t", plotter.BOXPLOT, fmt.Sprintf("plot type. Available types: %s", strings.Join([]string{plotter.BOXPLOT, plotter.LINE, plotter.NTH_SPEED}, ", ")))
	pflag.IntVarP(&index, "index", "i", 0, "position to plot the speeds in 'nth' charts")
	pflag.IntVarP(&clubID, "club", "c", 0, "club ID for which to load the data")
	pflag.IntVarP(&leagueID, "league", "l", 0, "league ID for which to load the data")
	pflag.IntVarP(&flagID, "flag", "f", 0, "flagID for which to load the data")
	pflag.StringVarP(&gender, "gender", "g", types.GENDER_MALE, "gender filter")
	pflag.StringVar(&category, "category", types.CATEGORY_ABSOLUT, "category filter")
	pflag.VarP(&years, "years", "y", "years to include in the data (can specify multiple times)")
	pflag.IntVarP(&day, "day", "d", 1, "day of the race for multiday races")
	pflag.BoolVar(&leaguesOnly, "leagues-only", false, "only races from a league")
	pflag.BoolVar(&branchTeams, "branch-teams", false, "filter only branch teams")
	pflag.BoolVarP(&normalize, "normalize", "n", false, "exclude outliers based on the speeds' standard deviation")
	pflag.StringVarP(&output, "output", "o", "", "saves the output plot")
	pflag.BoolVarP(&verbose, "verbose", "v", false, "define log level to DEBUG")

	s := service.Init()
	config := parseArgs(s)

	if verbose {
		prettylog.SetLevel(prettylog.DEBUG)
		prettylog.Debug("config=%v", *config)
	}

	plotter.PlotStats(s, config)
}

func parseArgs(service *service.Service) *plotter.PlotConfig {
	pflag.Parse()

	assert.Contains(gender, []any{types.GENDER_ALL, types.GENDER_MALE, types.GENDER_FEMALE, types.GENDER_MIX}, "invalid gender=%s", gender)
	assert.Contains(category, []any{types.CATEGORY_ABSOLUT, types.CATEGORY_SCHOOL, types.CATEGORY_VETERAN}, "invalid category=%s", category)
	assert.Contains(plotType, []any{plotter.BOXPLOT, plotter.LINE, plotter.NTH_SPEED}, "invalid plotType=%s", plotType)
	assert.Assert(plotType != plotter.NTH_SPEED || len(years) > 0, "plotType=%s requires at least one year", plotType)
	assert.Assert(plotType != plotter.NTH_SPEED || index > 0, "plotType=%s requires an index", plotType)

	validBoxplot := plotType == plotter.BOXPLOT && (clubID > 0 || leagueID > 0)
	validNthPlot := plotType == plotter.NTH_SPEED && leagueID > 0 && len(years) > 0 && index > 0
	validLinePlot := plotType == plotter.LINE && (clubID > 0 || flagID > 0) && len(years) > 0
	assert.Assert(validBoxplot || validNthPlot || validLinePlot, "invalid plot configuration")

	var err error
	var club *types.Entity
	if clubID > 0 {
		club, err = service.GetClubByID(int64(clubID))
		assert.NotNil(club, "invalid clubID=%d", clubID)
		assert.NoError(err, "invalid clubID=%d", clubID)
	}

	var flag *types.Flag
	if flagID > 0 {
		flag, err = service.GetFlagByID(int64(flagID))
		assert.NotNil(flag, "invalid flagID=%d", flagID)
		assert.NoError(err, "invalid flagID=%d", flagID)
	}

	var league *types.League
	if leagueID > 0 {
		league, err = service.GetLeagueByID(int64(leagueID))
		assert.NotNil(league, "invalid leagueID=%d", leagueID)
		assert.NoError(err, "invalid leagueID=%d", leagueID)

		if branchTeams {
			prettylog.Info("branch_teams is not supported with leagues, ignoring it")
			branchTeams = false
		}

		if league.Gender != nil && gender != *league.Gender {
			prettylog.Info("given gender=%s does not match %s, using league's one", gender, *league.Gender)
			gender = *league.Gender
		}

		if league.Category != nil && category != *league.Category {
			prettylog.Info("given category=%s does not match %s, using league's one", category, *league.Category)
			category = *league.Category
		}
	}

	return &plotter.PlotConfig{
		Index:       index,
		Club:        club,
		League:      league,
		Flag:        flag,
		PlotType:    plotType,
		Gender:      gender,
		Category:    category,
		Years:       years,
		Day:         day,
		Normalize:   normalize,
		LeaguesOnly: leaguesOnly,
		BranchTeams: branchTeams,
		Output:      output,
	}
}

var (
	plotType string
	index    int
	clubID   int
	leagueID int
	flagID   int
	gender   string
	category string
	years    yearsFlag
	day      int

	leaguesOnly bool
	branchTeams bool
	normalize   bool
	output      string
	verbose     bool
)

type yearsFlag []int

func (y *yearsFlag) String() string {
	return fmt.Sprint(*y)
}

func (y *yearsFlag) Set(value string) error {
	if strings.Contains(value, ",") {
		years := strings.Split(value, ",")
		for _, year := range years {
			var parsedYear int
			fmt.Sscanf(year, "%d", &parsedYear)
			*y = append(*y, parsedYear)
		}
	} else if strings.Contains(value, "..") {
		limits := strings.Split(value, "..")
		assert.Assert(len(limits) == 2, "invalid year limits=%s", limits)

		var start, end int
		fmt.Sscanf(limits[0], "%d", &start)
		fmt.Sscanf(limits[1], "%d", &end)
		assert.Assert(start > 0 && end > 0, "invalid year limits: start=%d end=%d", start, end)

		for i := start; i <= end; i++ {
			*y = append(*y, i)
		}
	} else {
		var year int
		fmt.Sscanf(value, "%d", &year)
		*y = append(*y, year)
	}
	return nil
}

func (y *yearsFlag) Type() string {
	return "yearsFlag"
}
