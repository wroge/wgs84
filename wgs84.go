package wgs84

import (
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
)

func LonLat() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid: spheroid.WGS84(),
		System:   system.LonLat{},
	}
}

func WebMercator() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid: spheroid.WGS84(),
		System:   system.LonLat{},
	}
}
