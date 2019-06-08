package osgb36

import (
	"github.com/wroge/wgs84"
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
	"github.com/wroge/wgs84/transformation"
)

func XYZ() wgs84.CoordinateReferenceSystem {
	return System(system.XYZ())
}

func LonLat() wgs84.CoordinateReferenceSystem {
	return System(system.LonLat())
}

func NationalGrid() wgs84.CoordinateReferenceSystem {
	return System(system.TransverseMercator(-2, 49, 0.9996012717, 400000, -100000))
}

func System(sys system.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.Airy(),
		Transformation: transformation.WGS84().Helmert(446.448, -125.157, 542.06, 0.15, 0.247, 0.842, -20.489),
		System:         sys,
	}
}
