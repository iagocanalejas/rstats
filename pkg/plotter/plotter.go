package plotter

import (
	"fmt"
	"math"
	"os/exec"
	"sort"
	"strconv"
	"sync"

	"github.com/iagocanalejas/rstats/internal/service"
	"github.com/iagocanalejas/rstats/internal/types"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
	prettylog "github.com/iagocanalejas/rstats/internal/utils/pretty-log"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	BOXPLOT   = "boxplot"
	LINE      = "line"
	NTH_SPEED = "nth"
)

type PlotConfig struct {
	Index  int
	Club   *types.Entity
	League *types.League
	Flag   *types.Flag

	PlotType string

	Gender   string
	Category string

	Years []int
	Day   int

	Normalize   bool
	LeaguesOnly bool
	BranchTeams bool

	Output string
}

func PlotStats(s *service.Service, config *PlotConfig) error {
	prettylog.Info("loading data")
	label := label(config.Index, config.Club, config.League, config.Normalize)

	var years *[]int
	var data *map[int][]float64
	var err error
	if config.PlotType == NTH_SPEED {
		years = &config.Years
		sort.Ints(*years)

		var wg sync.WaitGroup
		var mu sync.Mutex

		d := make(map[int][]float64)
		for _, year := range *years {
			wg.Add(1)

			go func(year int) {
				defer wg.Done()

				prettylog.Debug("loading data for year=%d", year)
				speeds, err := s.GetNthSpeedsBy(&service.GetNthSpeedsByParams{
					Index:           config.Index,
					Club:            config.Club,
					League:          config.League,
					Gender:          config.Gender,
					Category:        config.Category,
					Day:             int16(config.Day),
					Year:            int16(year),
					BranchTeams:     config.BranchTeams,
					OnlyLeagueRaces: config.LeaguesOnly,
					Normalize:       config.Normalize,
				})
				assert.NoError(err, "loading data for year=%d with config=%v", year, *config)

				mu.Lock()
				d[year] = *speeds
				mu.Unlock()
			}(year)
		}

		wg.Wait()
		data = &d
	} else {
		years, data, err = s.GetYearSpeedsBy(&service.GetYearSpeedsByParams{
			Club:            config.Club,
			League:          config.League,
			Flag:            config.Flag,
			Gender:          config.Gender,
			Category:        config.Category,
			Day:             int16(config.Day),
			Years:           config.Years,
			BranchTeams:     config.BranchTeams,
			OnlyLeagueRaces: config.LeaguesOnly,
			Normalize:       config.Normalize,
		})
		assert.NoError(err, "loading data with config=%v", *config)
	}

	switch config.PlotType {
	case BOXPLOT:
		return boxplot(label, data, years, config.Output)
	case LINE:
		return lineplot(label, data, years, config.Output)
	case NTH_SPEED:
		return lineplot(label, data, years, config.Output)
	}

	return nil
}

func boxplot(label string, data *map[int][]float64, years *[]int, output string) error {
	prettylog.Info("boxplotting")
	p := plot.New()

	boxplotIdx := 0
	for _, year := range *years {
		speeds := (*data)[year]
		values := make(plotter.Values, len(speeds))
		for i, v := range speeds {
			values[i] = v
		}

		boxplot, err := plotter.NewBoxPlot(vg.Points(20), float64(boxplotIdx), values)
		assert.NoError(err, "plotting boxplot", year, values)

		p.Add(boxplot)
		boxplotIdx++
	}

	p.Title.Text = label
	p.X.Label.Text = "Año"
	p.X.Tick.Marker = yearMarker{Years: *years}
	p.Y.Label.Text = "Velocidades"
	p.Y.Tick.Marker = quarterTicker{}

	err := displayOrSave(p, output)
	assert.NoError(err, "saving file")
	return err
}

func lineplot(label string, data *map[int][]float64, years *[]int, output string) error {
	prettylog.Info("lineplotting")
	p := plot.New()

	lines := make([]plot.Plotter, len(*years))
	lineplotIdx := 0
	maxValues := 0
	for _, year := range *years {
		pts := make(plotter.XYs, len((*data)[year]))
		for i := range (*data)[year] {
			pts[i].X = float64(i)
			pts[i].Y = (*data)[year][i]
		}
		line, err := plotter.NewLine(pts)
		assert.NoError(err, "error generating line", year)

		line.Color = plotutil.DefaultColors[lineplotIdx]
		lines[lineplotIdx] = line

		p.Legend.Add(strconv.Itoa(year), line)

		if len((*data)[year]) > maxValues {
			maxValues = len((*data)[year])
		}

		lineplotIdx++
	}

	p.Add(lines...)

	// X axis labels
	keys := make([]string, 0, maxValues)
	for k := range maxValues {
		if k%5 == 0 {
			keys = append(keys, strconv.Itoa(k))
		} else {
			keys = append(keys, "")
		}
	}

	p.Title.Text = label
	p.X.Label.Text = "Regata"
	p.X.Tick.Marker = keysMarker{Keys: keys}
	p.Y.Label.Text = "Velocidades"
	p.Y.Tick.Marker = quarterTicker{}

	err := displayOrSave(p, output)
	assert.NoError(err, "saving file")
	return err
}

func displayOrSave(p *plot.Plot, output string) error {
	if output == "" {
		filename := "./tmp/_temp_plot.png"

		err := p.Save(8*vg.Inch, 4*vg.Inch, filename)
		assert.NoError(err, "saving temporal file")

		prettylog.Info("opening plot")
		return exec.Command("xdg-open", filename).Start()
	}

	return p.Save(8*vg.Inch, 4*vg.Inch, output)
}

func label(index int, club *types.Entity, league *types.League, normalized bool) string {
	label := "VELOCIDADES (km/h)"
	if club != nil && league != nil {
		label = fmt.Sprintf("%s (%s) %s", club.Name, league.Symbol, label)
	} else if club != nil {
		label = fmt.Sprintf("%s %s", club.Name, label)
	} else if league != nil {
		label = fmt.Sprintf("%s %s", league.Symbol, label)
	}

	if index > 0 {
		label = fmt.Sprintf("VELOCIDADES (km/h) del %d", index)
		if league != nil {
			label = fmt.Sprintf("%s (%s)", label, league.Symbol)
		}
	}

	if normalized {
		label = fmt.Sprintf("%s - normalizadas", label)
	}

	return label
}

type quarterTicker struct{}

func (t quarterTicker) Ticks(minimum, maximum float64) []plot.Tick {
	var ticks []plot.Tick
	for i := math.Floor(minimum); i <= math.Ceil(maximum); i += 0.25 {
		if i == float64(int(i)) {
			ticks = append(ticks, plot.Tick{Value: i, Label: fmt.Sprintf("%.0f", i)})
		} else {
			ticks = append(ticks, plot.Tick{Value: i, Label: ""})
		}
	}
	return ticks
}

type yearMarker struct{ Years []int }

func (y yearMarker) Ticks(minimum, maximum float64) []plot.Tick {
	var ticks []plot.Tick
	for i := minimum; i <= maximum; i++ {
		year := y.Years[int(i)]
		if year%2 == 0 {
			ticks = append(ticks, plot.Tick{Value: float64(i), Label: fmt.Sprintf("%d", year)})
		} else {
			ticks = append(ticks, plot.Tick{Value: float64(i), Label: ""})
		}
	}
	return ticks
}

type keysMarker struct{ Keys []string }

func (k keysMarker) Ticks(minimum, maximum float64) []plot.Tick {
	var ticks []plot.Tick
	for i := minimum; i <= maximum; i++ {
		ticks = append(ticks, plot.Tick{Value: float64(i), Label: k.Keys[int(i)]})
	}
	return ticks
}
