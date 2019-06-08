package dhdn2001

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

func GK(zone float64) wgs84.CoordinateReferenceSystem {
	return System(system.GK(zone))
}

func System(sys system.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.Bessel(),
		Transformation: transformation.WGS84().Helmert(598.1, 73.7, 418.2, 0.202, 0.045, -2.455, 6.7),
		System:         sys,
	}
}
