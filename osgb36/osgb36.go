// Package osgb36 implements wgs84.CoordinateReferenceSystem's
// in the OSGB36 geodetic datum.
package osgb36

import (
	"github.com/wroge/wgs84"
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
	"github.com/wroge/wgs84/transformation"
)

// LonLat is equal to EPSG-Code 4277
func LonLat() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LonLat())
}

// NationalGrid is equal to EPSG-Code 27700
func NationalGrid() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.TransverseMercator(-2, 49, 0.9996012717, 400000, -100000))
}

// WithSystem provides any wgs84.System in the OSGB36
// geodetic datum.
func WithSystem(sys wgs84.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.Airy(),
		Transformation: transformation.WGS84().Helmert(446.448, -125.157, 542.06, 0.15, 0.247, 0.842, -20.489),
		System:         sys,
	}
}
