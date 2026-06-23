package wgs84

import (
	"math"
)

type WebMercator struct{}

func (w WebMercator) String() string {
	return "+proj=merc"
}

func (w WebMercator) ToGeographic(east, north, h float64, s Spheroid) (lon, lat, h2 float64) {
	D := (-north) / s.A
	phi := math.Pi/2 - 2*math.Atan(math.Pow(math.E, D))
	lambda := east / s.A

	return degree(lambda), degree(phi), h
}

func (w WebMercator) FromGeographic(lon, lat, h float64, s Spheroid) (east, north, h2 float64) {
	lambda := radian(lon)
	phi := radian(lat)

	east = s.A * lambda
	north = s.A * math.Log(math.Tan(math.Pi/4+phi/2))

	return east, north, h
}
