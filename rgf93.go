package wgs84

import (
	"github.com/wroge/wgs84/spheroid"
)

// EPSG-Code 2154
func RGF93FranceLambert() CoordinateReferenceSystem {
	return RGF93().LambertConformalConic2SP(3, 46.5, 49, 44, 700000, 6600000)
}

// CC42 - CC50, EPSG-Code 3942 - 3950
func RGF93CC(lat float64) CoordinateReferenceSystem {
	return RGF93().LambertConformalConic2SP(3, lat, lat-0.75, lat+0.75, 1700000, 2200000+(lat-43)*1000000)
}

func RGF93() GeodeticDatum {
	return GeodeticDatum{
		Spheroid: spheroid.GRS80(),
	}
}
