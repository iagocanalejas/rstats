package db

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/iagocanalejas/rstats/internal/utils"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
	prettylog "github.com/iagocanalejas/rstats/internal/utils/pretty-log"
	"github.com/lib/pq"
)

type Participant struct {
	ID int64 `db:"id"`

	Gender   string `db:"gender"`
	Category string `db:"category"`
	Distance int    `db:"distance"`

	ClubId       int64           `db:"club_id"`
	ClubName     string          `db:"club_name"`
	ClubRawNames *pq.StringArray `db:"club_raw_names"`

	IsDisqualified bool `db:"disqualified"`

	Laps   *pq.StringArray `db:"laps"`
	Lane   *int16          `db:"lane"`
	Series *int16          `db:"series"`
}

func (r *Repository) GetParticipantsByRaceID(raceID int64) ([]Participant, error) {
	query, args, err := sq.
		Select("p.id", "p.gender", "p.category", "p.distance", "p.laps", "p.lane", "p.series", "p.club_id as club_id", "e.name as club_name", "p.club_names as club_raw_names",
			"((SELECT count(*) FROM penalty pe WHERE pe.participant_id = p.id AND disqualification) > 0) as disqualified").
		From("participant p").
		LeftJoin("entity e ON p.club_id = e.id").
		Where(sq.Eq{"p.race_id": raceID}).
		OrderBy("p.laps[ARRAY_UPPER(p.laps, 1)] ASC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	assert.NoError(err, "building query", "query", query, "args", args)

	var participants []Participant
	if err = r.db.Select(&participants, query, args...); err != nil {
		return make([]Participant, 0), err
	}

	return participants, nil
}

type GetYearSpeedsByParams struct {
	ClubID          int64
	LeagueID        int64
	FlagID          int64
	Gender          string
	Category        string
	Day             int
	Years           []int
	BranchTeams     bool
	OnlyLeagueRaces bool
	Normalize       bool
}

// GetYearSpeedsBy retrieves the speeds of participants grouped by year based on the provided filtering criteria.
// It constructs and executes a SQL query that aggregates speeds for each year, applying optional normalization and year filtering.
//
// Params:
// - params (*GetYearSpeedsByParams): A struct containing the following fields to filter and control the query:
//   - ClubID (int): Filter races by the ID of the club.
//   - LeagueID (int): Filter races by the ID of the league.
//   - FlagID (int): Filter races by a specific flag (e.g., competition flag).
//   - Gender (string): Filter by participant gender.
//   - Category (string): Filter by race category.
//   - Day (time.Time): A specific day to filter the races.
//   - Years ([]int): A list of years to filter races.
//   - BranchTeams (bool): Whether to include branch teams in the query.
//   - OnlyLeagueRaces (bool): A flag to filter only league races.
//   - Normalize (bool): Whether to normalize speeds by excluding outliers.
//
// Returns:
// - (*[]int): A pointer to a slice of integers representing the years for which data was retrieved.
// - (*map[int][]float64): A pointer to a map where the key is the year, and the value is a slice of speeds for that year.
// - (error): An error if the query execution or row scanning fails.
//
// SQL Query Explanation:
//  1. **Speed Calculation**: Speed is calculated for each participant by dividing the race distance by the time taken,
//     then converting the result to km/h.
//  2. **Subquery**: Filters are applied to the races and participants based on the provided parameters (e.g., ClubID, LeagueID).
//  3. **Normalization** (optional): If enabled, speeds that fall outside two standard deviations from the mean are excluded.
//  4. **Main Query**: Aggregates speeds for each year using `array_agg`, and groups the results by year.
func (r *Repository) GetYearSpeedsBy(params *GetYearSpeedsByParams) (*[]int, *map[int][]float64, error) {
	subqueryWhere := getSpeedFilters(
		params.ClubID, params.LeagueID, params.FlagID,
		params.Gender, params.Category,
		params.Day,
		params.BranchTeams, params.OnlyLeagueRaces,
	)
	speedExpression := "(p.distance / (extract(EPOCH FROM p.laps[cardinality(p.laps)]))) * 3.6"

	whereClause := ""
	if len(params.Years) > 0 {
		whereClause = fmt.Sprintf("WHERE year in (%s)", utils.IntSlice2String(params.Years))
	}

	if params.Normalize {
		normalizeClause := `
			speed BETWEEN
			(
				SELECT AVG(speed) - (2 * STDDEV_POP(speed))
				FROM speeds_query
			)
			AND
			(
				SELECT AVG(speed) + (2 * STDDEV_POP(speed))
				FROM speeds_query
			)
		`
		if whereClause != "" {
			whereClause = whereClause + " AND " + normalizeClause
		} else {
			whereClause = "WHERE " + normalizeClause
		}
	}

	rawQuery := fmt.Sprintf(`
		WITH speeds_query AS (
            SELECT
                extract(YEAR from date)::INTEGER as year,
                CAST(%s AS DOUBLE PRECISION) as speed
            FROM participant p JOIN race r ON p.race_id = r.id
            WHERE %s
            ORDER BY r.date, speed DESC
        )
        SELECT year, array_agg(speed) AS speeds
        FROM speeds_query
        %s
        GROUP BY year
        ORDER BY year;
	`, speedExpression, subqueryWhere, whereClause)

	prettylog.Debug("%s", rawQuery)

	rows, err := r.db.Query(rawQuery)
	assert.NoError(err, "failed to execute query", "query", rawQuery)
	defer rows.Close()

	years := make([]int, 0)
	speeds := make(map[int][]float64)

	var year int
	var speedArray []float64

	for rows.Next() {
		err := rows.Scan(&year, pq.Array(&speedArray))
		assert.NoError(err, "failed to scan row", "year", year, "speeds", speedArray)

		years = append(years, year)
		speeds[year] = append([]float64(nil), speedArray...)
	}

	return &years, &speeds, nil
}

type GetNthSpeedsByParams struct {
	Index           int // the index is one-based as postgresql arrays are one-based
	ClubID          int64
	LeagueID        int64
	Gender          string
	Category        string
	Day             int
	Year            int
	BranchTeams     bool
	OnlyLeagueRaces bool
	Normalize       bool
}

// GetNthSpeedsBy retrieves the N-th highest speed for each race based on the provided filtering criteria.
// It constructs and executes a SQL query to extract speeds from race data, apply filters, and calculate the N-th
// highest speed using the provided index. Optionally, it can normalize the speeds by filtering out outliers.
//
// Params:
// - params (*GetNthSpeedsByParams): A struct containing the following fields used to filter and control the query:
//   - Index (int): The N-th position of the speed to retrieve (e.g., 1st fastest, 2nd fastest).
//   - ClubID (int): Filter races by the ID of the club.
//   - LeagueID (int): Filter races by the ID of the league.
//   - Gender (string): Filter by participant gender.
//   - Category (string): Filter by race category.
//   - Day (time.Time): The specific day to filter the races.
//   - Year (int): The year of the races.
//   - BranchTeams (bool): Whether to include branch teams in the query.
//   - OnlyLeagueRaces (bool): Flag to filter only league races.
//   - Normalize (bool): Whether to apply normalization, filtering out speeds more than two standard deviations from the mean.
//
// Returns:
// - (*[]float64): A pointer to a slice containing the N-th highest speeds for each race that meets the filtering criteria.
// - (error): An error if the query execution or row scanning fails.
//
// SQL Query Explanation:
//  1. **Speed Calculation**: The speed for each participant is calculated by dividing the race distance by the time taken
//     (in seconds) and converting it to km/h.
//  2. **Subquery**: Filters are applied to the races and participants based on the provided parameters (e.g., ClubID, Gender, Year).
//  3. **Normalization** (optional): If normalization is enabled, speeds outside two standard deviations from the mean are excluded.
//  4. **Main Query**: Retrieves the N-th highest speed for each race using `array_agg` and returns only races where there are at least N speeds.
func (r *Repository) GetNthSpeedsBy(params *GetNthSpeedsByParams) (*[]float64, error) {
	assert.Assert(params.Index > 0, "no index provided", *params)
	assert.Assert(params.Year > 0, "no year provided", *params)

	subqueryWhere := getSpeedFilters(
		params.ClubID, params.LeagueID, 0,
		params.Gender, params.Category,
		params.Day,
		params.BranchTeams, params.OnlyLeagueRaces,
	)
	subqueryWhere += fmt.Sprintf(" AND extract(YEAR FROM r.date) = %d", params.Year)
	speedExpression := "(p.distance / (extract(EPOCH FROM p.laps[cardinality(p.laps)]))) * 3.6"

	whereClause := ""
	if params.Normalize {
		whereClause = `
            WHERE speed BETWEEN (
                SELECT AVG(speed) - (2 * STDDEV_POP(speed))
                FROM speeds_query
            ) AND (
                SELECT AVG(speed) + (2 * STDDEV_POP(speed))
                FROM speeds_query
            )
		`
	}

	rawQuery := fmt.Sprintf(`
        WITH speeds_query AS (
            SELECT
                p.race_id,
                CAST(%s AS DOUBLE PRECISION) as speed
            FROM participant p JOIN race r ON p.race_id = r.id
            WHERE %s
            ORDER BY r.date
        )
        SELECT race_id, (array_agg(speed ORDER BY speed DESC))[%d] AS speed
        FROM speeds_query
		%s
        GROUP BY race_id
        HAVING array_length(array_agg(speed), 1) >= %d;
	`, speedExpression, subqueryWhere, params.Index, whereClause, params.Index)

	prettylog.Debug("%s", rawQuery)

	rows, err := r.db.Query(rawQuery)
	assert.NoError(err, "failed to execute query", "query", rawQuery)
	defer rows.Close()

	speeds := make([]float64, 0)

	var raceID int
	var speed float64
	for rows.Next() {
		err := rows.Scan(&raceID, &speed)
		assert.NoError(err, "failed to scan row", "raceID", raceID, "speed", speed)

		speeds = append(speeds, speed)
	}

	return &speeds, nil
}

func getSpeedFilters(
	clubID, leagueID, flagID int64,
	gender, category string,
	day int,
	branchTeams, onlyLeagueRaces bool,
) string {
	assert.Assert(gender != "", "invalid gender")
	assert.Assert(category != "", "invalid category")
	assert.Assert(day == 1 || day == 2, "invalid day")

	genderFilter := fmt.Sprintf("(p.gender = '%s' AND (r.gender = '%s' OR r.gender = '%s'))", gender, gender, "ALL")
	if onlyLeagueRaces || leagueID > 0 {
		genderFilter = fmt.Sprintf("(p.gender = '%s' AND r.gender = '%s')", gender, gender)
	}

	categoryFilter := fmt.Sprintf("(p.category = '%s' AND (r.category = '%s' OR r.category = '%s'))", category, category, "ALL")
	if onlyLeagueRaces || leagueID > 0 {
		categoryFilter = fmt.Sprintf("(p.category = '%s' AND r.category = '%s')", category, category)
	}

	filters := []string{
		"NOT r.cancelled",
		fmt.Sprintf("r.day = %d", day),
		"p.laps <> '{}'",
		"NOT p.retired",
		"NOT p.guest",
		"NOT p.absent",
		"p.distance IS NOT NULL",
		"(extract(EPOCH FROM p.laps[cardinality(p.laps)])) > 0",                              // Avoid division by zero
		"NOT EXISTS(SELECT * FROM penalty WHERE participant_id = p.id AND disqualification)", // Avoid disqualifications
		genderFilter,
		categoryFilter,
	}

	if branchTeams {
		filters = append(filters, "EXISTS(SELECT 1 FROM unnest(p.club_names) AS club_name WHERE club_name LIKE '% B')")
	} else if leagueID > 0 && flagID > 0 {
		filters = append(filters, "(p.club_names = '{}' OR NOT EXISTS(SELECT 1 FROM unnest(p.club_names) AS club_name WHERE club_name LIKE '% B'))")
	}

	if onlyLeagueRaces {
		filters = append(filters, "r.league_id IS NOT NULL")
	}

	if clubID > 0 {
		filters = append(filters, fmt.Sprintf("p.club_id = %d", clubID))
	}
	if leagueID > 0 {
		filters = append(filters, fmt.Sprintf("r.league_id = %d", leagueID))
	}
	if flagID > 0 {
		filters = append(filters, fmt.Sprintf("r.flag_id = %d", flagID))
	}

	return strings.Join(filters, " AND ")
}
