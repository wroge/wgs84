package system

import (
	"math"
)

type WebMercator struct{}

func (WebMercator) ToXYZ(east, north, h float64, s Spheroid) (x, y, z float64) {
	sph := getSpheroid(s)
	lon := degree(east / sph.A())
	lat := math.Atan(math.Exp(north/sph.A()))*degree(1)*2 - 90
	return LonLat{}.ToXYZ(lon, lat, h, sph)
}

func (WebMercator) FromXYZ(x, y, z float64, s Spheroid) (east, north, h float64) {
	sph := getSpheroid(s)
	lon, lat, h := LonLat{}.FromXYZ(x, y, z, sph)
	east = radian(lon) * sph.A()
	north = math.Log(math.Tan(radian((90+lat)/2))) * sph.A()
	return
}
