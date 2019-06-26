// Package gda94 implements wgs84.CoordinateReferenceSystem's
// in the GDA94 geodetic datum.
package gda94

import (
	"github.com/wroge/wgs84"
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
	"github.com/wroge/wgs84/transformation"
)

// LonLat is equal to EPSG-Code 4283.
func LonLat() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LonLat())
}

// CKIG94 is equal to EPSG-Code 6723.
func CKIG94() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.TransverseMercator(96.875, 0, 0.99999387, 50000, 1500000))
}

// CIG94 is equal to EPSG-Code 6721
func CIG94() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.TransverseMercator(105.625, 0, 1.00002514, 50000, 1300000))
}

// BCSG02 is equal to EPSG-Code 3113
func BCSG02() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.TransverseMercator(153, -28, 0.99999, 50000, 100000))
}

// Vicgrid94 is equal to EPSG-Code 3111
func Vicgrid94() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(145, -37, -36, -38, 2500000, 2500000))
}

// AustralianAlbers is equal to EPSG-Code 3577
func AustralianAlbers() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.AlbersEqualAreaConic(132, 0, -18, -36, 0, 0))
}

// NSWLambert is equal to EPSG-Code 3308
func NSWLambert() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(147, -33.25, -30.75, -35.75, 9300000, 4500000))
}

// SALambert is equal to EPSG-Code 3107
func SALambert() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(135, -32, -28, -36, 1000000, 2000000))
}

// UTM provides UTM south coordinates in the GDA94 geodetic datum.
//
// UTM Zone 49 is equal to EPSG-Code 28349,
// UTM Zone 50 is equal to EPSG-Code 28350.
func UTM(zone float64) wgs84.CoordinateReferenceSystem {
	return WithSystem(system.UTM(zone, false))
}

// WithSystem provides any wgs84.System in the GDA94
// geodetic datum.
func WithSystem(sys wgs84.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.GRS80(),
		Transformation: transformation.WGS84(),
		System:         sys,
	}
}
