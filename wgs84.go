package wgs84

import (
	"github.com/wroge/wgs84/spheroid"
)

// EPSG-Code 4326
func LonLat() CoordinateReferenceSystem {
	return WGS84().LonLat()
}

// EPSG-Code 3857, 900913
func WebMercator() CoordinateReferenceSystem {
	return WGS84().WebMercator()
}

func WGS84() GeodeticDatum {
	return GeodeticDatum{
		Spheroid: spheroid.WGS84(),
	}
}
