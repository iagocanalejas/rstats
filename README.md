# Plot Data

Different type of plots to visualize stats.

```sh
go run cmd/plot/main.go \
	[-t, --type TYPE] \
	[-i, --index INDEX] \
	[-c, --club CLUB_ID] \
	[-l, --league LEAGUE_ID] \
	[-f, --flag FLAG_ID] \
	[-g, --gender GENDER] \
	[--category CATEGORY] \
	[-y, --years YEARS] \
	[-d, --day DAY] \
	[-n, --normalize] \
	[--leagues-only] \
	[--branch-teams] \
	[-o, --output FILE] \
	[-v, --verbose]

# options:
#   -t TYPE, --type TYPE
#                         plot type ['boxplot', 'line', 'nth'].
#   -i INDEX, --index INDEX
#                         position to plot the speeds in 'nth' charts.
#   -c CLUB, --club CLUB
#                         club ID for which to load the data.
#   -l LEAGUE, --league LEAGUE
#                         league ID for which to load the data.
#   -f FLAG, --flag FLAG
#                         flag ID for which to load the data.
#   -g GENDER, --gender GENDER
#                         gender filter.
#   --category CATEGORY
#                         category filter.
#   -y [YEARS ...], --years [YEARS ...]
#                         years to include in the data.
#   -d DAY, --day DAY
#                         day of the race for multiday races.
#   -o OUTPUT, --output OUTPUT
#                         saves the output plot.
#   --leagues-only
#                         only races from a league.
#   --branch-teams
#                         filter only branch teams.
#   -n, --normalize
#                         exclude outliers based on the speeds' standard deviation.
#   -v, --verbose
#                         increase output verbosity.
```

### Examples

```sh
# Plot the winner speed of each race for the league 5 in 2015, 2016, 2017, and 2018.
# The plot will be saved in the Downloads folder with the name test.png.
go run cmd/plot/main.go -t nth --league 5 -i 1 -y 2015..2018 -o ~/Downloads/p.png
```

```sh
# Plot the normalized league speeds of the Puebla team for all the years.
go run cmd/plot/main.go -c 25 --leagues-only -n -o ~/Downloads/p.png
```

```sh
# Plot all AVG speeds per year for the league 5.
go run cmd/plot/main.go --league 5 -o ~/Downloads/p.png
```

```sh
# Plot the speeds of the Puebla team for the league 5 in 2021, 2022, and 2023.
go run cmd/plot/main.go -t line -c 25 --league 5 -y 2021..2023 -o ~/Downloads/p.png
```

```sh
# Plot the speeds of the Puebla team for the flag 12 in 2021, 2022, and 2023.
go run cmd/plot/main.go -t line -f 12 -y 2021..2023 -o ~/Downloads/p.png
```

```sh
# Plot all leagues AVG speeds per year.
# The plot will be saved in the Downloads folder with a generated name.
parallel -j 11 go run cmd/plot/main.go --league {} -o ~/Downloads/l{}.png ::: $(seq 1 11)
```

# Search Outliers

Search for outliers in the data.

```sh
go run cmd/outliers/main.go \
	[-t, --threshold, THRESHOLD] \
	[--exclude EXCLUDED_RACE_IDS] \
	[--limit, -l] \
	[-v, --verbose]

# options:
#   -t THRESHOLD, --threshold THRESHOLD
#                         threshold to consider a value an outlier.
#   --exclude EXCLUDED_RACE_IDS
#                         races that will be ignored.
#   -l, --limit
#                         apply spped limits (8.0, 20.0) for outlier detection.
#   -v, --verbose
#                         increase output verbosity.
```

# Boat Calculator

Given boat data, calculate some of the floating properties.

- Center of Gravity (CG)
- Compute ratio of flotability

```sh
go run cmd/boat/main.go FILE_NAME
```

# Terminal UI for regatas

```sh
# to run the TUI use
go run cmd/tui/main.go
```

# Technologies

## PostgreSQL

## [SQLX](https://github.com/jmoiron/sqlx)

## [Squirrel](https://github.com/Masterminds/squirrel)

## [TView](https://github.com/rivo/tview)
