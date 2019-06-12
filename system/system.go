// Package system implements the wgs84.System interface.
package system

import (
	"math"

	"github.com/wroge/wgs84/spheroid"
)

// The unnamed Spheroid interface is implemented by
// github.com/wroge/wgs84/spheroid spheroid.Spheroid.
type Spheroid = interface {
	A() float64
	Fi() float64
}

// System implements the wgs84.System interface.
type System struct {
	toXYZ   func(a, b, c float64, s Spheroid) (x, y, z float64)
	fromXYZ func(x, y, z float64, s Spheroid) (a, b, c float64)
}

// ToXYZ is used in the wgs84.System interface.
func (sys System) ToXYZ(a, b, c float64, s Spheroid) (x, y, z float64) {
	if sys.toXYZ == nil {
		return LonLat().toXYZ(a, b, c, s)
	}
	return sys.toXYZ(a, b, c, s)
}

// FromXYZ is used in the wgs84.System interface.
func (sys System) FromXYZ(x, y, z float64, s Spheroid) (a, b, c float64) {
	if sys.fromXYZ == nil {
		return LonLat().fromXYZ(x, y, z, s)
	}
	return sys.fromXYZ(x, y, z, s)
}

// XYZ is a geocentric coordinate System.
func XYZ() System {
	return System{
		toXYZ: func(a, b, c float64, s Spheroid) (x, y, z float64) {
			return a, b, c
		},
		fromXYZ: func(x, y, z float64, s Spheroid) (a, b, c float64) {
			return x, y, z
		},
	}
}

// LonLat is a geocentric coordinate System.
func LonLat() System {
	N := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() / math.Sqrt(1-sph.E2()*math.Pow(math.Sin(φ), 2))
	}
	return System{
		toXYZ: func(lon, lat, h float64, s Spheroid) (x, y, z float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			x = (N(radian(lat), sph) + h) * math.Cos(radian(lon)) * math.Cos(radian(lat))
			y = (N(radian(lat), sph) + h) * math.Cos(radian(lat)) * math.Sin(radian(lon))
			z = (N(radian(lat), sph)*math.Pow(sph.A()*(1-sph.F()), 2)/(sph.A2()) + h) * math.Sin(radian(lat))
			return
		},
		fromXYZ: func(x, y, z float64, s Spheroid) (lon, lat, h float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			sd := math.Sqrt(x*x + y*y)
			T := math.Atan(z * sph.A() / (sd * sph.B()))
			B := math.Atan((z + sph.E2()*(sph.A2())/sph.B()*math.Pow(math.Sin(T), 3)) / (sd - sph.E2()*sph.A()*math.Pow(math.Cos(T), 3)))
			h = sd/math.Cos(B) - N(B, sph)
			return degree(math.Atan2(y, x)), degree(B), h
		},
	}
}

// TransverseMercator is a projected coordinate System.
// False Longitude, False Latitude, 1. Parrallel Latitude, 2. Parrallel Latitude,
// False Easting, False Northing.
func TransverseMercator(lonf, latf, scale, eastf, northf float64) System {
	M := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() * ((1-sph.E2()/4-3*sph.E4()/64-5*sph.E6()/256)*φ -
			(3*sph.E2()/8+3*sph.E4()/32+45*sph.E6()/1024)*math.Sin(2*φ) +
			(15*sph.E4()/256+45*sph.E6()/1024)*math.Sin(4*φ) -
			(35*sph.E6()/3072)*math.Sin(6*φ))
	}
	N := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() / math.Sqrt(1-sph.E2()*sin2(φ))
	}
	T := func(φ float64) float64 {
		return tan2(φ)
	}
	C := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.Ei2() * cos2(φ)
	}
	return System{
		toXYZ: func(east, north, h float64, s Spheroid) (x, y, z float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			east -= eastf
			north -= northf
			Mi := M(radian(latf), sph) + north/scale
			μ := Mi / (sph.A() * (1 - sph.E2()/4 - 3*sph.E4()/64 - 5*sph.E6()/256))
			φ1 := μ + (3*sph.Ei()/2-27*sph.Ei3()/32)*math.Sin(2*μ) +
				(21*sph.Ei2()/16-55*sph.Ei4()/32)*math.Sin(4*μ) +
				(151*sph.Ei3()/96)*math.Sin(6*μ) +
				(1097*sph.Ei4()/512)*math.Sin(8*μ)
			R1 := sph.A() * (1 - sph.E2()) / math.Pow(1-sph.E2()*sin2(φ1), 3/2)
			D := east / (N(φ1, sph) * scale)
			φ := φ1 - (N(φ1, sph)*math.Tan(φ1)/R1)*(D*D/2-(5+3*T(φ1)+10*C(φ1, sph)-4*C(φ1, sph)*C(φ1, sph)-9*sph.Ei2())*
				math.Pow(D, 4)/24+(61+90*T(φ1)+298*C(φ1, sph)+45*T(φ1)*T(φ1)-252*sph.Ei2()-3*C(φ1, sph)*C(φ1, sph))*
				math.Pow(D, 6)/720)
			λ := radian(lonf) + (D-(1+2*T(φ1)+C(φ1, sph))*D*D*D/6+(5-2*C(φ1, sph)+28*T(φ1)-3*C(φ1, sph)*C(φ1, sph)+8*sph.Ei2()+24*T(φ1)*T(φ1))*
				math.Pow(D, 5)/120)/math.Cos(φ1)
			return LonLat().toXYZ(degree(λ), degree(φ), h, sph)
		},
		fromXYZ: func(x, y, z float64, s Spheroid) (east, north, h float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			lon, lat, h := LonLat().fromXYZ(x, y, z, sph)
			φ := radian(lat)
			A := (radian(lon) - radian(lonf)) * math.Cos(φ)
			east = scale*N(φ, sph)*(A+(1-T(φ)+C(φ, sph))*
				math.Pow(A, 3)/6+(5-18*T(φ)+T(φ)*T(φ)+72*C(φ, sph)-58*sph.Ei2())*math.Pow(A, 5)/120) + eastf
			north = scale*(M(φ, sph)-M(radian(latf), sph)+N(φ, sph)*math.Tan(φ)*(A*A/2+(5-T(φ)+9*C(φ, sph)+4*C(φ, sph)*C(φ, sph))*
				math.Pow(A, 4)/24+(61-58*T(φ)+T(φ)*T(φ)+600*C(φ, sph)-330*sph.Ei2())*math.Pow(A, 6)/720)) + northf
			return
		},
	}
}

// UTM is a projected coordinate System (TransverseMercator).
// It has 6 degrees zone width and a scale factor of 0.9996
// at the central meridian.
func UTM(zone float64, northern bool) System {
	if northern {
		return TransverseMercator(zone*6-183, 0, 0.9996, 500000, 0)
	}
	return TransverseMercator(zone*6-183, 0, 0.9996, 500000, 10000000)
}

// GK is a projected coordinate System (TransverseMercator).
// It has 3 degrees zone width.
func GK(zone float64) System {
	return TransverseMercator(zone*3, 0, 1, zone*1000000+500000, 0)
}

// Mercator is a projected coordinate System. The parameters:
// False Longitude, Scale Factor, False Easting, False Northing.
func Mercator(lonf, scale, eastf, northf float64) System {
	return System{
		toXYZ: func(east, north, h float64, s Spheroid) (x, y, z float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			east = (east - eastf) / scale
			north = (north - northf) / scale
			t := math.Exp(-north * sph.A())
			φ := math.Pi/2 - 2*math.Atan(t)
			for i := 0; i < 5; i++ {
				φ = math.Pi/2 - 2*math.Atan(t*math.Pow((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ)), sph.E()/2))
			}
			return LonLat().toXYZ(east/sph.A()+lonf, degree(φ), h, sph)
		},
		fromXYZ: func(x, y, z float64, s Spheroid) (east, north, h float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			lon, lat, h := LonLat().fromXYZ(x, y, z, sph)
			east = scale * sph.A() * (radian(lon) - radian(lonf))
			north = scale * sph.A() / 2 *
				math.Log(1+math.Sin(radian(lat))/(1-math.Sin(radian(lat)))*math.Pow((1-sph.E()*math.Sin(radian(lat)))/(1+sph.E()*math.Sin(radian(lat))), math.E))
			return
		},
	}
}

// WebMercator is a projected coordinate System.
func WebMercator() System {
	return System{
		toXYZ: func(east, north, h float64, s Spheroid) (x, y, z float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			lon := degree(east / s.A())
			lat := math.Atan(math.Exp(north/s.A()))*degree(1)*2 - 90
			return LonLat().toXYZ(lon, lat, h, s)
		},
		fromXYZ: func(x, y, z float64, s Spheroid) (east, north, h float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			lon, lat, h := LonLat().fromXYZ(x, y, z, s)
			east = radian(lon) * s.A()
			north = math.Log(math.Tan(radian((90+lat)/2))) * s.A()
			return
		},
	}
}

// LambertConformalConic1SP is a projected coordinate System. The parameters:
// False Longitude, False Latitude, Scale Factor, False Easting, False Northing.
func LambertConformalConic1SP(lonf, latf, scale, eastf, northf float64) System {
	t := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Tan(math.Pi/4-φ/2) /
			math.Pow((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ)), sph.E()/2)
	}
	m := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Cos(φ) / math.Sqrt(1-sph.E2()*sin2(φ))
	}
	n := math.Sin(radian(latf))
	F := func(sph spheroid.Spheroid) float64 {
		return m(radian(latf), sph) / (n * math.Pow(t(radian(latf), sph), n))
	}
	ρ := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() * F(sph) * math.Pow(t(φ, sph)*scale, n)
	}
	return System{
		toXYZ: func(east, north, h float64, s Spheroid) (x, y, z float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			ρi := math.Sqrt(math.Pow(east-eastf, 2) + math.Pow(ρ(radian(latf), sph)-(north-northf), 2))
			if n < 0 {
				ρi = -ρi
			}
			ti := math.Pow(ρi/(sph.A()*scale*F(sph)), 1/n)
			φ := math.Pi/2 - 2*math.Atan(ti)
			for i := 0; i < 5; i++ {
				φ = math.Pi/2 - 2*math.Atan(ti*math.Pow((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ)), sph.E()/2))
			}
			λ := math.Atan((east-eastf)/(ρ(radian(latf), sph)-(north-northf)))/n + radian(lonf)
			return LonLat().toXYZ(degree(λ), degree(φ), h, sph)
		},
		fromXYZ: func(x, y, z float64, s Spheroid) (east, north, h float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			lon, lat, h := LonLat().fromXYZ(x, y, z, sph)
			θ := n * (radian(lon) - radian(lonf))
			east = eastf + ρ(radian(lat), sph)*math.Sin(θ)
			north = northf + ρ(radian(latf), sph) - ρ(radian(lat), sph)*math.Cos(θ)
			return
		},
	}
}

// LambertConformalConic2SP is a projected coordinate System. The parameters:
// False Longitude, False Latitude, 1. Parrallel Latitude, 2. Parrallel Latitude,
// False Easting, False Northing.
func LambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf float64) System {
	t := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Tan(math.Pi/4-φ/2) /
			math.Pow((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ)), sph.E()/2)
	}
	m := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Cos(φ) / math.Sqrt(1-sph.E2()*sin2(φ))
	}
	n := func(sph spheroid.Spheroid) float64 {
		if radian(lat1) == radian(lat2) {
			return math.Sin(radian(lat1))
		}
		return (math.Log(m(radian(lat1), sph)) - math.Log(m(radian(lat2), sph))) / (math.Log(t(radian(lat1), sph)) - math.Log(t(radian(lat2), sph)))
	}
	F := func(sph spheroid.Spheroid) float64 {
		return m(radian(lat1), sph) / (n(sph) * math.Pow(t(radian(lat1), sph), n(sph)))
	}
	ρ := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() * F(sph) * math.Pow(t(φ, sph), n(sph))
	}
	return System{
		toXYZ: func(east, north, h float64, s Spheroid) (x, y, z float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			ρi := math.Sqrt(math.Pow(east-eastf, 2) + math.Pow(ρ(radian(latf), sph)-(north-northf), 2))
			if n(sph) < 0 {
				ρi = -ρi
			}
			ti := math.Pow(ρi/(sph.A()*F(sph)), 1/n(sph))
			φ := math.Pi/2 - 2*math.Atan(ti)
			for i := 0; i < 5; i++ {
				φ = math.Pi/2 - 2*math.Atan(ti*math.Pow((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ)), sph.E()/2))
			}
			λ := math.Atan((east-eastf)/(ρ(radian(latf), sph)-(north-northf)))/n(sph) + radian(lonf)
			return LonLat().toXYZ(degree(λ), degree(φ), h, sph)
		},
		fromXYZ: func(x, y, z float64, s Spheroid) (east, north, h float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			lon, lat, h := LonLat().fromXYZ(x, y, z, sph)
			θ := n(sph) * (radian(lon) - radian(lonf))
			east = eastf + ρ(radian(lat), sph)*math.Sin(θ)
			north = northf + ρ(radian(latf), sph) - ρ(radian(lat), sph)*math.Cos(θ)
			return
		},
	}
}

// AlbersEqualAreaConic is a projected coordinate System. The parameters:
// False Longitude, False Latitude, 1. Parrallel Latitude, 2. Parrallel Latitude,
// False Easting, False Northing.
func AlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf float64) System {
	m := func(φ float64, sph spheroid.Spheroid) float64 {
		return math.Cos(φ) / math.Sqrt(1-sph.E2()*sin2(φ))
	}
	q := func(φ float64, sph spheroid.Spheroid) float64 {
		return (1 - sph.E2()) * (math.Sin(φ)/(1-sph.E2()*sin2(φ)) -
			(1/(2*sph.E()))*math.Log((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ))))
	}
	n := func(sph spheroid.Spheroid) float64 {
		if radian(lat1) == radian(lat2) {
			return math.Sin(radian(lat1))
		}
		return (m(radian(lat1), sph)*m(radian(lat1), sph) - m(radian(lat2), sph)*m(radian(lat2), sph)) / (q(radian(lat2), sph) - q(radian(lat1), sph))
	}
	C := func(sph spheroid.Spheroid) float64 {
		return m(radian(lat1), sph)*m(radian(lat1), sph) + n(sph)*q(radian(lat1), sph)
	}
	ρ := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() * math.Sqrt(C(sph)-n(sph)*q(φ, sph)) / n(sph)
	}
	return System{
		toXYZ: func(east, north, h float64, s Spheroid) (x, y, z float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			east -= eastf
			north -= northf
			ρi := math.Sqrt(east*east + math.Pow(ρ(radian(latf), sph)-north, 2))
			qi := (C(sph) - ρi*ρi*n(sph)*n(sph)/sph.A2()) / n(sph)
			φ := math.Asin(qi / 2)
			for i := 0; i < 5; i++ {
				φ += math.Pow(1-sph.E2()*sin2(φ), 2) /
					(2 * math.Cos(φ)) * (qi/(1-sph.E2()) -
					math.Sin(φ)/(1-sph.E2()*sin2(φ)) +
					1/(2*sph.E())*math.Log((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ))))
			}
			θ := math.Atan(east / (ρ(radian(latf), sph) - north))
			return LonLat().toXYZ(degree(radian(lonf)+θ/n(sph)), degree(φ), h, sph)
		},
		fromXYZ: func(x, y, z float64, s Spheroid) (east, north, h float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			lon, lat, h := LonLat().fromXYZ(x, y, z, sph)
			θ := n(sph) * (radian(lon) - radian(lonf))
			east = eastf + ρ(radian(lat), sph)*math.Sin(θ)
			north = northf + ρ(radian(latf), sph) - ρ(radian(lat), sph)*math.Cos(θ)
			return
		},
	}
}

// EquidistantConic is a projected coordinate System. The parameters:
// False Longitude, False Latitude, 1. Parrallel Latitude, 2. Parrallel Latitude,
// False Easting, False Northing.
func EquidistantConic(lonf, latf, lat1, lat2, eastf, northf float64) System {
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
		if radian(lat1) == radian(lat2) {
			return math.Sin(radian(lat1))
		}
		return sph.A() * (m(radian(lat1), sph) - m(radian(lat2), sph)) / (M(radian(lat2), sph) - M(radian(lat1), sph))
	}
	G := func(sph spheroid.Spheroid) float64 {
		return m(radian(lat1), sph)/n(sph) + M(radian(lat1), sph)/sph.A()
	}
	ρ := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A()*G(sph) - M(φ, sph)
	}
	return System{
		toXYZ: func(east, north, h float64, s Spheroid) (x, y, z float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			east -= eastf
			north -= northf
			ρi := math.Sqrt(east*east + math.Pow(ρ(radian(latf), sph)-north, 2))
			if n(sph) < 0 {
				ρi = -ρi
			}
			Mi := sph.A()*G(sph) - ρi
			μ := Mi / (sph.A() * (1 - sph.E2()/4 - 3*sph.E4()/64 - 5*sph.E6()/256))
			φ := μ + (3*sph.Ei()/2-27*sph.Ei3()/32)*math.Sin(2*μ) +
				(21*sph.Ei2()/16-55*sph.Ei4()/32)*math.Sin(4*μ) +
				(151*sph.Ei3()/96)*math.Sin(6*μ) +
				(1097*sph.Ei4()/512)*math.Sin(8*μ)
			θ := math.Atan(east / (ρ(radian(latf), sph) - north))
			return LonLat().toXYZ(degree((radian(lonf) + θ/n(sph))), degree(φ), h, sph)
		},
		fromXYZ: func(x, y, z float64, s Spheroid) (east, north, h float64) {
			if s == nil {
				s = spheroid.WGS84()
			}
			sph := spheroid.New(s.A(), s.Fi())
			lon, lat, h := LonLat().fromXYZ(x, y, z, sph)
			θ := n(sph) * (radian(lon) - radian(lonf))
			east = eastf + ρ(radian(lat), sph)*math.Sin(θ)
			north = northf + ρ(radian(latf), sph) - ρ(radian(lat), sph)*math.Cos(θ)
			return
		},
	}
}

func sin2(east float64) float64 {
	return math.Pow(math.Sin(east), 2)
}

func cos2(east float64) float64 {
	return math.Pow(math.Cos(east), 2)
}

func tan2(east float64) float64 {
	return math.Pow(math.Tan(east), 2)
}

func degree(r float64) float64 {
	return r * 180 / math.Pi
}

func radian(d float64) float64 {
	return d * math.Pi / 180
}
