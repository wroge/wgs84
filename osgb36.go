package wgs84

import (
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/transformation"
)

func OSGB36() GeodeticDatum {
	return GeodeticDatum{
		Spheroid:       spheroid.Airy(),
		Transformation: transformation.Helmert{446.448, -125.157, 542.06, 0.15, 0.247, 0.842, -20.489},
	}
}

func OSGB36LonLat() CoordinateReferenceSystem {
	return OSGB36().LonLat()
}

func OSGB36NationalGrid() CoordinateReferenceSystem {
	return OSGB36().TransverseMercator(-2, 49, 0.9996012717, 400000, -100000)
}
