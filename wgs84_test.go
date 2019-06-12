package wgs84

import (
	"fmt"

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
	east, north, h := EPSG27700.FromLonLat()(-1, 52, 0)
	fmt.Println(east, north, h)
}

func ExampleTo() {
	// import (
	// 	"github.com/wroge/wgs84/etrs89"
	// )
	EPSG25832 := CoordinateReferenceSystem{
		Spheroid:       spheroid.GRS80(),
		Transformation: nil,
		System:         system.UTM(32, true),
	}
	east, north, h := To(EPSG25832)(9, 52, 0)
	fmt.Println(east, north, h)
}

func ExampleCoordinateReferenceSystem_ToWebMercator() {
	east, north, h := LonLat().ToWebMercator().Round(0)(9, 52, 0)
	fmt.Println(east, north, h)
}
