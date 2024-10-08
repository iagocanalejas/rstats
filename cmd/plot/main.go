package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/iagocanalejas/rstats/internal/db"
	"github.com/iagocanalejas/rstats/internal/service"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	setupFileLogger()

	clubID, err := strconv.Atoi(os.Args[1])
	if err != nil || clubID == 0 {
		panic("clubID is required")
	}

	filters := &db.ParticipantByClubFilters{}

	var is_female bool
	flag.BoolVar(&is_female, "female", false, "female races.")
	if is_female {
		filters.Gender = "FEMALE"
	} else {
		filters.Gender = "MALE"
	}

	flag.BoolVar(&filters.OnlyLeagueRaces, "league", false, "only races from a league.")
	flag.BoolVar(&filters.BranchTeams, "branch", false, "only branch teams.")

	service := service.Init()
	speeds, err := service.GetSpeedAVGByClubID(int64(clubID), filters)
	if err != nil {
		panic(err)
	}

	p := plot.New()

	years := make([]string, len(speeds))
	for i, speed := range speeds {
		values := make(plotter.Values, len(speed.Speeds))
		for i, v := range speed.Speeds {
			values[i] = v
		}

		boxplot, err := plotter.NewBoxPlot(vg.Points(20), float64(i), values)
		if err != nil {
			panic(err)
		}

		p.Add(boxplot)
		years[i] = strconv.Itoa(int(speed.Year))
	}

	p.Title.Text = "Speeds by Year"
	p.X.Label.Text = "Year"
	p.Y.Label.Text = "Speed"
	p.NominalX(years...)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "speeds_by_year.png"); err != nil {
		panic(err)
	}
}

func setupFileLogger() {
	f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	log.SetOutput(f)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
