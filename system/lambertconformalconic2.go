package system

import (
	"math"

	"github.com/wroge/wgs84/spheroid"
)

type LambertConformalConic2SP struct {
	Lonf, Latf, Lat1, Lat2, Eastf, Northf float64
}

func (sys LambertConformalConic2SP) ToXYZ(east, north, h float64, s Spheroid) (x, y, z float64) {
	t := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Tan(math.Pi/4-φ/2) /
			math.Pow((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ)), sph.E()/2)
	}
	m := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Cos(φ) / math.Sqrt(1-sph.E2()*sin2(φ))
	}
	n := func(sph spheroid.Spheroid) float64 {
		if radian(sys.Lat1) == radian(sys.Lat2) {
			return math.Sin(radian(sys.Lat1))
		}
		return (math.Log(m(radian(sys.Lat1), sph)) - math.Log(m(radian(sys.Lat2), sph))) /
			(math.Log(t(radian(sys.Lat1), sph)) - math.Log(t(radian(sys.Lat2), sph)))
	}
	F := func(sph spheroid.Spheroid) float64 {
		return m(radian(sys.Lat1), sph) / (n(sph) * math.Pow(t(radian(sys.Lat1), sph), n(sph)))
	}
	ρ := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() * F(sph) * math.Pow(t(φ, sph), n(sph))
	}
	sph := getSpheroid(s)
	ρi := math.Sqrt(math.Pow(east-sys.Eastf, 2) + math.Pow(ρ(radian(sys.Latf), sph)-(north-sys.Northf), 2))
	if n(sph) < 0 {
		ρi = -ρi
	}
	ti := math.Pow(ρi/(sph.A()*F(sph)), 1/n(sph))
	φ := math.Pi/2 - 2*math.Atan(ti)
	for i := 0; i < 5; i++ {
		φ = math.Pi/2 - 2*math.Atan(ti*math.Pow((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ)), sph.E()/2))
	}
	λ := math.Atan((east-sys.Eastf)/(ρ(radian(sys.Latf), sph)-(north-sys.Northf)))/n(sph) + radian(sys.Lonf)
	return LonLat{}.ToXYZ(degree(λ), degree(φ), h, sph)
}

func (sys LambertConformalConic2SP) FromXYZ(x, y, z float64, s Spheroid) (east, north, h float64) {
	t := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Tan(math.Pi/4-φ/2) /
			math.Pow((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ)), sph.E()/2)
	}
	m := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Cos(φ) / math.Sqrt(1-sph.E2()*sin2(φ))
	}
	n := func(sph spheroid.Spheroid) float64 {
		if radian(sys.Lat1) == radian(sys.Lat2) {
			return math.Sin(radian(sys.Lat1))
		}
		return (math.Log(m(radian(sys.Lat1), sph)) - math.Log(m(radian(sys.Lat2), sph))) /
			(math.Log(t(radian(sys.Lat1), sph)) - math.Log(t(radian(sys.Lat2), sph)))
	}
	F := func(sph spheroid.Spheroid) float64 {
		return m(radian(sys.Lat1), sph) / (n(sph) * math.Pow(t(radian(sys.Lat1), sph), n(sph)))
	}
	ρ := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() * F(sph) * math.Pow(t(φ, sph), n(sph))
	}
	sph := getSpheroid(s)
	lon, lat, h := LonLat{}.FromXYZ(x, y, z, sph)
	θ := n(sph) * (radian(lon) - radian(sys.Lonf))
	east = sys.Eastf + ρ(radian(lat), sph)*math.Sin(θ)
	north = sys.Northf + ρ(radian(sys.Latf), sph) - ρ(radian(lat), sph)*math.Cos(θ)
	return
}
