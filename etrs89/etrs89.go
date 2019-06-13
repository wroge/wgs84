package etrs89

import (
	"github.com/wroge/wgs84"
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
	"github.com/wroge/wgs84/transformation"
)

// UTM is a universal transverse mercator system in the
// ETRS89 geodetic datum.
//
// UTM(32) is equal to EPSG-Code 25832
// UTM(33) is equal to EPSG-Code 25833 ...
func UTM(zone float64) wgs84.CoordinateReferenceSystem {
	return WithSystem(system.UTM(zone, true))
}

// WithSystem provides any wgs84.System in the ETRS89
// geodetic datum.
func WithSystem(sys wgs84.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.GRS80(),
		Transformation: transformation.WGS84(),
		System:         sys,
	}
}
