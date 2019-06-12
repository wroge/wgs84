package wgs84

import (
	"github.com/wroge/wgs84/etrs89"
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
	"github.com/wroge/wgs84/transformation"
)

func ExampleCoordinateReferenceSystem() {
	var EPSG27700 = CoordinateReferenceSystem{
		Spheroid:       spheroid.Airy(),
		Transformation: transformation.WGS84().Helmert(446.448, -125.157, 542.06, 0.15, 0.247, 0.842, -20.489),
		System:         system.TransverseMercator(-2, 49, 0.9996012717, 400000, -100000),
	}
}

func ExampleTo() {
	// import (
	// 	"github.com/wroge/wgs84/etrs89"
	// )
	east, north, h := To(etrs89.UTM(32))(9, 52, 0)
}

func ExampleCoordinateReferenceSystem_ToWebMercator() {
	east, north, h := LonLat().ToWebMercator().Round(0)(9, 52, 0)
}
