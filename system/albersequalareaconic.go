package system

import (
	"math"

	"github.com/wroge/wgs84/spheroid"
)

type AlbersEqualAreaConic struct {
	Lonf, Latf, Lat1, Lat2, Eastf, Northf float64
}

func (sys AlbersEqualAreaConic) ToXYZ(east, north, h float64, s Spheroid) (x, y, z float64) {
	sph := getSpheroid(s)
	east -= sys.Eastf
	north -= sys.Northf
	ρi := math.Sqrt(east*east + math.Pow(sys._ρ(radian(sys.Latf), sph)-north, 2))
	qi := (sys._C(sph) - ρi*ρi*sys._n(sph)*sys._n(sph)/sph.A2()) / sys._n(sph)
	φ := math.Asin(qi / 2)
	for i := 0; i < 5; i++ {
		φ += math.Pow(1-sph.E2()*sin2(φ), 2) /
			(2 * math.Cos(φ)) * (qi/(1-sph.E2()) -
			math.Sin(φ)/(1-sph.E2()*sin2(φ)) +
			1/(2*sph.E())*math.Log((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ))))
	}
	θ := math.Atan(east / (sys._ρ(radian(sys.Latf), sph) - north))
	return LonLat{}.ToXYZ(degree(radian(sys.Lonf)+θ/sys._n(sph)), degree(φ), h, sph)
}

func (sys AlbersEqualAreaConic) FromXYZ(x, y, z float64, s Spheroid) (east, north, h float64) {
	sph := getSpheroid(s)
	lon, lat, h := LonLat{}.FromXYZ(x, y, z, sph)
	θ := sys._n(sph) * (radian(lon) - radian(sys.Lonf))
	east = sys.Eastf + sys._ρ(radian(lat), sph)*math.Sin(θ)
	north = sys.Northf + sys._ρ(radian(sys.Latf), sph) - sys._ρ(radian(lat), sph)*math.Cos(θ)
	return
}

func (AlbersEqualAreaConic) _m(φ float64, sph spheroid.Spheroid) float64 {
	return math.Cos(φ) / math.Sqrt(1-sph.E2()*sin2(φ))
}

func (AlbersEqualAreaConic) _q(φ float64, sph spheroid.Spheroid) float64 {
	return (1 - sph.E2()) * (math.Sin(φ)/(1-sph.E2()*sin2(φ)) -
		(1/(2*sph.E()))*math.Log((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ))))
}

func (sys AlbersEqualAreaConic) _n(sph spheroid.Spheroid) float64 {
	if radian(sys.Lat1) == radian(sys.Lat2) {
		return math.Sin(radian(sys.Lat1))
	}
	return (sys._m(radian(sys.Lat1), sph)*sys._m(radian(sys.Lat1), sph) - sys._m(radian(sys.Lat2), sph)*sys._m(radian(sys.Lat2), sph)) /
		(sys._q(radian(sys.Lat2), sph) - sys._q(radian(sys.Lat1), sph))
}

func (sys AlbersEqualAreaConic) _C(sph spheroid.Spheroid) float64 {
	return sys._m(radian(sys.Lat1), sph)*sys._m(radian(sys.Lat1), sph) + sys._n(sph)*sys._q(radian(sys.Lat1), sph)
}

func (sys AlbersEqualAreaConic) _ρ(φ float64, sph spheroid.Spheroid) float64 {
	return sph.A() * math.Sqrt(sys._C(sph)-sys._n(sph)*sys._q(φ, sph)) / sys._n(sph)
}
