package wgs84

import (
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
)

// EPSG-Code 4326
func LonLat() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid: spheroid.WGS84(),
		System:   system.LonLat{},
	}
}

// EPSG-Code 3857, 900913
func WebMercator() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid: spheroid.WGS84(),
		System:   system.LonLat{},
	}
}
