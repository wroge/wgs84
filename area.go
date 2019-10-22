package wgs84

import "math"

type AreaFunc func(lon, lat float64) bool

func (a AreaFunc) Contains(lon, lat float64) bool {
	if math.Abs(lat) > 180 || math.Abs(lat) > 90 {
		return false
	}
	if a == nil {
		return true
	}
	return a(lon, lat)
}
