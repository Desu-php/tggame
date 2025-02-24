package utils

import "math"

func GrowthIncrease(currentValue float64, growthFactor float64) float64 {
	increase := currentValue * (growthFactor / 100)
	return math.Round(currentValue + increase)
}
