package wgs84

import (
	"github.com/wroge/wgs84/spheroid"
)

func ETRS89() GeodeticDatum {
	return GeodeticDatum{
		Spheroid: spheroid.GRS80(),
	}
}

func ETRS89UTM(zone float64) CoordinateReferenceSystem {
	return ETRS89().UTM(zone, true)
}
