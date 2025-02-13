package service

import "github.com/iagocanalejas/rstats/internal/types"

func (s *Service) ComputeCenterOfGravity(config *types.BoatConfig) float64 {
	totalWeight := 0.0
	totalMoment := 0.0
	for i, seat := range config.Seats {
		totalWeight += seat.Weight
		totalMoment += seat.Weight * s.computeWeightPositionFromBow(config, i)
	}
	return totalMoment / totalWeight
}

func (s *Service) ComputeRatio(config *types.BoatConfig) float64 {
	sum := 0.0
	for i := 0; i < len(config.Seats); i++ {
		sum += ((float64(config.Length) / 2) - s.computeWeightPositionFromBow(config, i)) * config.Seats[i].Weight
	}
	return sum / 1000
}

func (s *Service) computeWeightPositionFromBow(config *types.BoatConfig, position int) float64 {
	if position == 0 {
		return float64(config.Length - config.Seats[0].BenchDistance)
	}
	rowlockPosition := config.Rowlocks[position-1] - int(config.Seats[position].RowlockPosition*config.RowlockGap)
	return float64(rowlockPosition - config.Seats[position].BenchDistance)
}
