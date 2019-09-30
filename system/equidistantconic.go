package system

import (
	"math"

	"github.com/wroge/wgs84/spheroid"
)

type EquidistantConic struct {
	Lonf, Latf, Lat1, Lat2, Eastf, Northf float64
}

func (sys EquidistantConic) ToXYZ(east, north, h float64, s Spheroid) (x, y, z float64) {
	M := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() * ((1-sph.E2()/4-3*sph.E4()/64-5*sph.E6()/256)*φ -
			(3*sph.E2()/8+3*sph.E4()/32+45*sph.E6()/1024)*math.Sin(2*φ) +
			(15*sph.E4()/256+45*sph.E6()/1024)*math.Sin(4*φ) -
			(35*sph.E6()/3072)*math.Sin(6*φ))
	}
	m := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Cos(φ) / math.Sqrt(1-sph.E2()*sin2(φ))
	}
	n := func(sph spheroid.Spheroid) float64 {
		if radian(sys.Lat1) == radian(sys.Lat2) {
			return math.Sin(radian(sys.Lat1))
		}
		return sph.A() * (m(radian(sys.Lat1), sph) - m(radian(sys.Lat2), sph)) / (M(radian(sys.Lat2), sph) - M(radian(sys.Lat1), sph))
	}
	G := func(sph spheroid.Spheroid) float64 {
		return m(radian(sys.Lat1), sph)/n(sph) + M(radian(sys.Lat1), sph)/sph.A()
	}
	ρ := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A()*G(sph) - M(φ, sph)
	}
	sph := getSpheroid(s)
	east -= sys.Eastf
	north -= sys.Northf
	ρi := math.Sqrt(east*east + math.Pow(ρ(radian(sys.Latf), sph)-north, 2))
	if n(sph) < 0 {
		ρi = -ρi
	}
	Mi := sph.A()*G(sph) - ρi
	μ := Mi / (sph.A() * (1 - sph.E2()/4 - 3*sph.E4()/64 - 5*sph.E6()/256))
	φ := μ + (3*sph.Ei()/2-27*sph.Ei3()/32)*math.Sin(2*μ) +
		(21*sph.Ei2()/16-55*sph.Ei4()/32)*math.Sin(4*μ) +
		(151*sph.Ei3()/96)*math.Sin(6*μ) +
		(1097*sph.Ei4()/512)*math.Sin(8*μ)
	θ := math.Atan(east / (ρ(radian(sys.Latf), sph) - north))
	return LonLat{}.ToXYZ(degree((radian(sys.Lonf) + θ/n(sph))), degree(φ), h, sph)
}

func (sys EquidistantConic) FromXYZ(x, y, z float64, s Spheroid) (east, north, h float64) {
	M := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() * ((1-sph.E2()/4-3*sph.E4()/64-5*sph.E6()/256)*φ -
			(3*sph.E2()/8+3*sph.E4()/32+45*sph.E6()/1024)*math.Sin(2*φ) +
			(15*sph.E4()/256+45*sph.E6()/1024)*math.Sin(4*φ) -
			(35*sph.E6()/3072)*math.Sin(6*φ))
	}
	m := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Cos(φ) / math.Sqrt(1-sph.E2()*sin2(φ))
	}
	n := func(sph spheroid.Spheroid) float64 {
		if radian(sys.Lat1) == radian(sys.Lat2) {
			return math.Sin(radian(sys.Lat1))
		}
		return sph.A() * (m(radian(sys.Lat1), sph) - m(radian(sys.Lat2), sph)) / (M(radian(sys.Lat2), sph) - M(radian(sys.Lat1), sph))
	}
	G := func(sph spheroid.Spheroid) float64 {
		return m(radian(sys.Lat1), sph)/n(sph) + M(radian(sys.Lat1), sph)/sph.A()
	}
	ρ := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A()*G(sph) - M(φ, sph)
	}
	sph := getSpheroid(s)
	lon, lat, h := LonLat{}.FromXYZ(x, y, z, sph)
	θ := n(sph) * (radian(lon) - radian(sys.Lonf))
	east = sys.Eastf + ρ(radian(lat), sph)*math.Sin(θ)
	north = sys.Northf + ρ(radian(sys.Latf), sph) - ρ(radian(lat), sph)*math.Cos(θ)
	return
}
