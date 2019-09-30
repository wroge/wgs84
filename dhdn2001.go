package wgs84

import (
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/transformation"
)

func DHDN2001() GeodeticDatum {
	return GeodeticDatum{
		Spheroid:       spheroid.Bessel(),
		Transformation: transformation.Helmert{598.1, 73.7, 418.2, 0.202, 0.045, -2.455, 6.7},
	}
}

func DHDN2001GK(zone float64) CoordinateReferenceSystem {
	return DHDN2001().GK(zone)
}
