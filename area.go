package wgs84

import "math"

// AreaFunc provides the Area interface.
type AreaFunc func(lon, lat float64) bool

// Contains is the implementation of the Area interface.
//
// By default Contains is true between -180 to 180 degrees longitude and
// -90 and 90 degrees latitude.
func (a AreaFunc) Contains(lon, lat float64) bool {
	if math.Abs(lat) > 180 || math.Abs(lat) > 90 {
		return false
	}
	if a == nil {
		return true
	}
	return a(lon, lat)
}
