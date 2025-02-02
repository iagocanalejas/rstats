package main

import (
	"sync"

	"github.com/iagocanalejas/rstats/internal/service"
	"github.com/iagocanalejas/rstats/internal/types"
	"github.com/iagocanalejas/rstats/internal/utils/arrays"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
	prettylog "github.com/iagocanalejas/rstats/internal/utils/pretty-log"
	"github.com/spf13/pflag"
)

type ParticipantOutlier struct {
	RaceID        int64
	ParticipantID int64
	AVGSpeed      float64
	Speed         float64
}

func main() {
	pflag.Float64VarP(&threshold, "threshold", "t", 0.3, "threshold for outlier detection")
	pflag.Int64SliceVar(&excludedRaceIDs, "exclude", []int64{}, "define ignored race IDs")
	pflag.BoolVarP(&limits, "limits", "l", false, "apply spped limits (8.0, 20.0) for outlier detection")
	pflag.BoolVarP(&verbose, "verbose", "v", false, "define log level to DEBUG")
	pflag.Parse()

	if verbose {
		prettylog.SetLevel(prettylog.DEBUG)
		prettylog.Debug("threshold=%f, excludedRaceIDs=%v", threshold, excludedRaceIDs)
	}

	s := service.Init()

	participants, err := s.GetParticipantsWithSpeed()
	assert.NoError(err, "loading participants with speed")
	prettylog.Info("grouped into %d", len(participants))

	batchSize := 500
	ProcessOutliers(participants, batchSize)

	// batchSizes := []int{100, 500, 1000, 5000}
	// for _, batchSize := range batchSizes {
	// 	start := time.Now()
	// 	ProcessOutliers(participants, batchSize)
	// 	elapsed := time.Since(start)
	// 	fmt.Printf("Batch size: %d, Time taken: %s\n", batchSize, elapsed)
	// }
}

func ProcessOutliers(participants [][]*types.Participant, batchSize int) {
	outliers := make([]*ParticipantOutlier, 0)

	outliersChan := make(chan *ParticipantOutlier)
	doneChan := make(chan bool)
	go func() {
		for {
			outlier, ok := <-outliersChan
			if !ok {
				// ensure we close the doneChan when the channel is done processing
				doneChan <- true
				break
			}
			prettylog.Info("found outlier: %+v", *outlier)
			outliers = append(outliers, outlier)
		}
		prettylog.Info("finished processing all outliers from channel")
	}()

	var wg sync.WaitGroup
	for i := 0; i < len(participants); i += batchSize {
		j := i + batchSize
		if j > len(participants) {
			j = len(participants)
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			for _, group := range participants[start:end] {
				detectRaceOutliers(outliersChan, group, threshold)
			}
		}(i, j)
	}

	wg.Wait()
	close(outliersChan)
	<-doneChan // wait for the channel to finish processing

	prettylog.Info("found %d outliers", len(outliers))

	races := make([]int64, 0)
	for _, outlier := range outliers {
		if !arrays.Contains(races, outlier.RaceID) {
			races = append(races, outlier.RaceID)
		}
	}
	prettylog.Info("grouped in %d races", len(races))
	prettylog.Info("\nraces: %+v", races)
}

func detectRaceOutliers(channel chan *ParticipantOutlier, group []*types.Participant, threshold float64) {
	sum := 0.0
	for _, participant := range group {
		sum += *participant.Speed
	}
	mean := sum / float64(len(group))
	for _, participant := range group {
		if arrays.Contains(excludedRaceIDs, participant.RaceID) {
			return
		}
		if *participant.Speed > mean*(threshold+1) || *participant.Speed < mean*(1-threshold) {
			channel <- &ParticipantOutlier{
				RaceID:        participant.RaceID,
				ParticipantID: participant.ID,
				AVGSpeed:      mean,
				Speed:         *participant.Speed,
			}
			continue
		}
		if limits && (*participant.Speed > 20.0 || *participant.Speed < 8.0) {
			channel <- &ParticipantOutlier{
				RaceID:        participant.RaceID,
				ParticipantID: participant.ID,
				AVGSpeed:      0.0,
				Speed:         *participant.Speed,
			}
			continue
		}
	}
}

var (
	threshold       float64
	excludedRaceIDs []int64
	limits          bool
	verbose         bool
)
