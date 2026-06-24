package wgs84

import (
	"math"
)

func degree(r float64) float64 {
	return r * 180 / math.Pi
}

func radian(g float64) float64 {
	return g * math.Pi / 180
}

func intPow(val float64, times int) float64 {
	result := 1.0

	for range times {
		result *= val
	}

	return result
}

func sin2(r float64) float64 {
	return intPow(math.Sin(r), 2)
}

func sign(x float64) float64 {
	if x < 0 {
		return -1
	}
	return 1
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
