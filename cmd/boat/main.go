package main

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/iagocanalejas/rstats/internal/service"
	"github.com/iagocanalejas/rstats/internal/types"
	"github.com/iagocanalejas/rstats/internal/utils/assert"
	prettylog "github.com/iagocanalejas/rstats/internal/utils/pretty-log"
)

func main() {
	prettylog.SetLevel(prettylog.DEBUG)

	assert.Assert(len(os.Args) > 1, "no arguments provided")
	fileName := os.Args[1]

	assert.Assert(strings.HasSuffix(fileName, ".json"), "invalid file extension=%s", fileName)

	file, err := os.Open(fileName)
	assert.NoError(err, "error opening file=%s", fileName)
	defer file.Close()

	data, err := io.ReadAll(file)
	assert.NoError(err, "error reading file=%s", fileName)

	boatConfig := &types.BoatConfig{}
	err = json.Unmarshal(data, boatConfig)
	assert.NoError(err, "error unmarshalling data=%s", data)

	validateConfig(boatConfig)
	prettylog.Info("config=%+v", *boatConfig)

	s := service.Static()

	centerOfGravity := s.ComputeCenterOfGravity(boatConfig)
	prettylog.Debug("centerOfGravity=%v (from bow)", centerOfGravity)

	ratio := s.ComputeRatio(boatConfig)
	prettylog.Debug("ratio=%v (more means more weight in bow)", ratio)
}

func validateConfig(config *types.BoatConfig) {
	assert.Assert(config.Weight > 0, "invalid weight=%v", config.Weight)
	assert.Assert(config.Length > 0, "invalid length=%v", config.Length)
	assert.Assert(len(config.Rowlocks) == len(config.Seats)-1, "invalid rowlocks=%v", config.Rowlocks)
	assert.Assert(config.RowlockGap > 0, "invalid rowlockGap=%v", config.RowlockGap)
	// assert.Assert(len(config.Ribs) == len(config.Seats), "invalid ribs=%v", config.Ribs)
	assert.Assert(config.BowFloatingLine > 0, "invalid bowFloatingLine=%v", config.BowFloatingLine)
	assert.Assert(config.SternFloatingLine > 0, "invalid sternFloatingLine=%v", config.SternFloatingLine)

	for i, seat := range config.Seats {
		assert.Contains(seat.Position, []types.Position{types.COXSWAIN, types.STROKE, types.TWO, types.THREE, types.FOUR, types.FIVE, types.SIX, types.BOW}, "invalid position=%s for seat=%d", seat.Position, i)
		if seat.Side != nil {
			assert.Assert(seat.RowlockPosition > 0, "invalid rowlockPosition=%v for seat=%d", seat.RowlockPosition, i)
			assert.Contains(*seat.Side, []types.Side{types.STARBOARD, types.PORT}, "invalid side=%s for seat=%d", *seat.Side, i)
		}
		assert.Assert(seat.Weight > 0, "invalid weight=%v for seat=%d", seat.Weight, i)
		assert.Assert(seat.BenchDistance > 0, "invalid benchDistance=%v for seat=%d", seat.BenchDistance, i)
	}
}
