// Package dhdn2001 implements wgs84.CoordinateReferenceSystem's
// in the DHDN2001 geodetic datum.
package dhdn2001

import (
	"github.com/wroge/wgs84"
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
	"github.com/wroge/wgs84/transformation"
)

// GK is a Gauss Krueger system in the DHDN2001
// geodetic datum.
//
// GK(2) is equal to EPSG-Code 31466.
// GK(3) is equal to EPSG-Code 31467.
// GK(4) is equal to EPSG-Code 31468.
// GK(5) is equal to EPSG-Code 31469.
func GK(zone float64) wgs84.CoordinateReferenceSystem {
	return WithSystem(system.GK(zone))
}

// WithSystem provides any wgs84.System in the DHDN2001
// geodetic datum.
func WithSystem(sys wgs84.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.Bessel(),
		Transformation: transformation.WGS84().Helmert(598.1, 73.7, 418.2, 0.202, 0.045, -2.455, 6.7),
		System:         sys,
	}
}
