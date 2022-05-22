package wgs84

import "math"

// AreaFunc implements the Contains method of the Area interface.
//
// Returns true if nil.
//
// Returns false for latitudes with an absolute over 90 and longitudes over 180.
type AreaFunc func(lon, lat float64) bool

// Contains method is the implementation of the Area interface.
//
// Returns false for latitudes with an absolute over 90 and longitudes over 180.
func (a AreaFunc) Contains(lon, lat float64) bool {
	return math.Abs(lat) <= 180 && math.Abs(lat) <= 90 && (a == nil || a(lon, lat))
}
