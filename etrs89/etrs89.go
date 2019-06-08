package etrs89

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

func UTM(zone float64) wgs84.CoordinateReferenceSystem {
	return System(system.UTM(zone, true))
}

func System(sys system.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.GRS80(),
		Transformation: transformation.WGS84(),
		System:         sys,
	}
}
