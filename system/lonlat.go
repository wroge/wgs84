package system

import (
	"math"

	"github.com/wroge/wgs84/spheroid"
)

type LonLat struct{}

func (LonLat) ToXYZ(lon, lat, h float64, s Spheroid) (x, y, z float64) {
	N := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() / math.Sqrt(1-sph.E2()*math.Pow(math.Sin(φ), 2))
	}
	sph := getSpheroid(s)
	x = (N(radian(lat), sph) + h) * math.Cos(radian(lon)) * math.Cos(radian(lat))
	y = (N(radian(lat), sph) + h) * math.Cos(radian(lat)) * math.Sin(radian(lon))
	z = (N(radian(lat), sph)*math.Pow(sph.A()*(1-sph.F()), 2)/(sph.A2()) + h) * math.Sin(radian(lat))
	return
}

func (LonLat) FromXYZ(x, y, z float64, s Spheroid) (lon, lat, h float64) {
	N := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() / math.Sqrt(1-sph.E2()*math.Pow(math.Sin(φ), 2))
	}
	sph := getSpheroid(s)
	sd := math.Sqrt(x*x + y*y)
	T := math.Atan(z * sph.A() / (sd * sph.B()))
	B := math.Atan((z + sph.E2()*(sph.A2())/sph.B()*
		math.Pow(math.Sin(T), 3)) / (sd - sph.E2()*sph.A()*math.Pow(math.Cos(T), 3)))
	h = sd/math.Cos(B) - N(B, sph)
	return degree(math.Atan2(y, x)), degree(B), h
}
