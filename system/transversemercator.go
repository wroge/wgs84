package system

import (
	"math"

	"github.com/wroge/wgs84/spheroid"
)

func UTM(zone float64, northern bool) TransverseMercator {
	if northern {
		return TransverseMercator{zone*6 - 183, 0, 0.9996, 500000, 0}
	}
	return TransverseMercator{zone*6 - 183, 0, 0.9996, 500000, 10000000}
}

func GK(zone float64) TransverseMercator {
	return TransverseMercator{zone * 3, 0, 1, zone*1000000 + 500000, 0}
}

type TransverseMercator struct {
	Lonf, Latf, Scale, Eastf, Northf float64
}

func (sys TransverseMercator) ToXYZ(east, north, h float64, s Spheroid) (x, y, z float64) {
	M := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() * ((1-sph.E2()/4-3*sph.E4()/64-5*sph.E6()/256)*φ -
			(3*sph.E2()/8+3*sph.E4()/32+45*sph.E6()/1024)*math.Sin(2*φ) +
			(15*sph.E4()/256+45*sph.E6()/1024)*math.Sin(4*φ) -
			(35*sph.E6()/3072)*math.Sin(6*φ))
	}
	N := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() / math.Sqrt(1-sph.E2()*sin2(φ))
	}
	T := tan2
	C := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.Ei2() * cos2(φ)
	}
	sph := getSpheroid(s)
	east -= sys.Eastf
	north -= sys.Northf
	Mi := M(radian(sys.Latf), sph) + north/sys.Scale
	μ := Mi / (sph.A() * (1 - sph.E2()/4 - 3*sph.E4()/64 - 5*sph.E6()/256))
	φ1 := μ + (3*sph.Ei()/2-27*sph.Ei3()/32)*math.Sin(2*μ) +
		(21*sph.Ei2()/16-55*sph.Ei4()/32)*math.Sin(4*μ) +
		(151*sph.Ei3()/96)*math.Sin(6*μ) +
		(1097*sph.Ei4()/512)*math.Sin(8*μ)
	R1 := sph.A() * (1 - sph.E2()) / math.Pow(1-sph.E2()*sin2(φ1), 3/2)
	D := east / (N(φ1, sph) * sys.Scale)
	φ := φ1 - (N(φ1, sph)*math.Tan(φ1)/R1)*(D*D/2-(5+3*T(φ1)+10*C(φ1, sph)-4*C(φ1, sph)*C(φ1, sph)-9*sph.Ei2())*
		math.Pow(D, 4)/24+(61+90*T(φ1)+298*C(φ1, sph)+45*T(φ1)*T(φ1)-252*sph.Ei2()-3*C(φ1, sph)*C(φ1, sph))*
		math.Pow(D, 6)/720)
	λ := radian(sys.Lonf) + (D-(1+2*T(φ1)+C(φ1, sph))*D*D*D/6+(5-2*C(φ1, sph)+
		28*T(φ1)-3*C(φ1, sph)*C(φ1, sph)+8*sph.Ei2()+24*T(φ1)*T(φ1))*
		math.Pow(D, 5)/120)/math.Cos(φ1)
	return LonLat{}.ToXYZ(degree(λ), degree(φ), h, sph)
}

func (sys TransverseMercator) FromXYZ(x, y, z float64, s Spheroid) (east, north, h float64) {
	M := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() * ((1-sph.E2()/4-3*sph.E4()/64-5*sph.E6()/256)*φ -
			(3*sph.E2()/8+3*sph.E4()/32+45*sph.E6()/1024)*math.Sin(2*φ) +
			(15*sph.E4()/256+45*sph.E6()/1024)*math.Sin(4*φ) -
			(35*sph.E6()/3072)*math.Sin(6*φ))
	}
	N := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.A() / math.Sqrt(1-sph.E2()*sin2(φ))
	}
	T := tan2
	C := func(φ float64, sph spheroid.Spheroid) float64 {
		return sph.Ei2() * cos2(φ)
	}
	sph := getSpheroid(s)
	lon, lat, h := LonLat{}.FromXYZ(x, y, z, sph)
	φ := radian(lat)
	A := (radian(lon) - radian(sys.Lonf)) * math.Cos(φ)
	east = sys.Scale*N(φ, sph)*(A+(1-T(φ)+C(φ, sph))*
		math.Pow(A, 3)/6+(5-18*T(φ)+T(φ)*T(φ)+72*C(φ, sph)-58*sph.Ei2())*
		math.Pow(A, 5)/120) + sys.Eastf
	north = sys.Scale*(M(φ, sph)-M(radian(sys.Latf), sph)+N(φ, sph)*math.Tan(φ)*
		(A*A/2+(5-T(φ)+9*C(φ, sph)+4*C(φ, sph)*C(φ, sph))*
			math.Pow(A, 4)/24+(61-58*T(φ)+T(φ)*T(φ)+600*C(φ, sph)-330*sph.Ei2())*math.Pow(A, 6)/720)) + sys.Northf
	return
}
