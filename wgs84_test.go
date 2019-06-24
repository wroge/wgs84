package wgs84

import (
	"fmt"
	"testing"

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
	conversion := EPSG27700.FromLonLat()
	if conversion != nil {
		fmt.Println(conversion(-1, 52, 0))
	}
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
	transformation := To(EPSG25832)
	if transformation != nil {
		fmt.Println(transformation(9, 52, 0))
	}
}

func ExampleCoordinateReferenceSystem_ToWebMercator() {
	conversion := LonLat().ToWebMercator().Round(0)
	if conversion != nil {
		fmt.Println(conversion(9, 52, 0))
	}
}

func BenchmarkToWebMercator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToWebMercator()(9, 52, 0)
	}
}

func BenchmarkTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		To(WebMercator())(9, 52, 0)
	}
}

func BenchmarkCoordinateReferenceSystem_ToWebMercator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LonLat().ToWebMercator()(9, 52, 0)
	}
}
