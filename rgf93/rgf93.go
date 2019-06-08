package rgf93

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

func FranceLambert() wgs84.CoordinateReferenceSystem {
	return System(system.LambertConformalConic2SP(3, 46.5, 49, 44, 700000, 6600000))
}

func CC(latitude float64) wgs84.CoordinateReferenceSystem {
	return System(system.LambertConformalConic2SP(3, latitude, latitude-0.75, latitude+0.75, 1700000, 2200000+(latitude-43)*1000000))
}

func System(sys system.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.GRS80(),
		Transformation: transformation.WGS84(),
		System:         sys,
	}
}
