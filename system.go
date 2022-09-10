//nolint:varnamelen,nonamedreturns,gomnd,asciicheck,nosnakecase
package wgs84

import (
	"math"
)

type webMercator struct{}

func (p webMercator) ToLonLat(east, north float64, s Spheroid) (lon, lat float64) {
	sph := spheroid{a: s.A(), fi: s.Fi()}
	lon = degree(east / sph.A())
	lat = math.Atan(math.Exp(north/sph.A()))*degree(1)*2 - 90

	return lon, lat
}

func (p webMercator) FromLonLat(lon, lat float64, s Spheroid) (east, north float64) {
	sph := spheroid{a: s.A(), fi: s.Fi()}
	east = radian(lon) * sph.A()
	north = math.Log(math.Tan(radian((90+lat)/2))) * sph.A()

	return east, north
}

type transverseMercator struct {
	lonf, latf, scale, eastf, northf float64
}

func (p transverseMercator) ToLonLat(east, north float64, s Spheroid) (lon, lat float64) {
	sph := spheroid{a: s.A(), fi: s.Fi()}
	east -= p.eastf
	north -= p.northf
	Mi := p._M(radian(p.latf), sph) + north/p.scale
	μ := Mi / (sph.A() * (1 - sph.e2()/4 - 3*sph.e4()/64 - 5*sph.e6()/256))
	φ1 := μ + (3*sph.ei()/2-27*sph.ei3()/32)*math.Sin(2*μ) +
		(21*sph.ei2()/16-55*sph.ei4()/32)*math.Sin(4*μ) +
		(151*sph.ei3()/96)*math.Sin(6*μ) +
		(1097*sph.ei4()/512)*math.Sin(8*μ)
	R1 := sph.A() * (1 - sph.e2()) / math.Pow(1-sph.e2()*sin2(φ1), 3/2)
	D := east / (p._N(φ1, sph) * p.scale)
	φ := φ1 - (p._N(φ1, sph)*math.Tan(φ1)/R1)*(D*D/2-(5+3*p._T(φ1)+10*
		p._C(φ1, sph)-4*p._C(φ1, sph)*p._C(φ1, sph)-9*sph.ei2())*
		math.Pow(D, 4)/24+(61+90*p._T(φ1)+298*p._C(φ1, sph)+45*p._T(φ1)*
		p._T(φ1)-252*sph.ei2()-3*p._C(φ1, sph)*p._C(φ1, sph))*
		math.Pow(D, 6)/720)
	λ := radian(p.lonf) + (D-(1+2*p._T(φ1)+p._C(φ1, sph))*D*D*D/6+(5-2*p._C(φ1, sph)+
		28*p._T(φ1)-3*p._C(φ1, sph)*p._C(φ1, sph)+8*sph.ei2()+24*p._T(φ1)*p._T(φ1))*
		math.Pow(D, 5)/120)/math.Cos(φ1)

	return degree(λ), degree(φ)
}

func (p transverseMercator) FromLonLat(lon, lat float64, s Spheroid) (east, north float64) {
	sph := spheroid{a: s.A(), fi: s.Fi()}
	φ := radian(lat)
	A := (radian(lon) - radian(p.lonf)) * math.Cos(φ)
	east = p.scale*p._N(φ, sph)*(A+(1-p._T(φ)+p._C(φ, sph))*
		math.Pow(A, 3)/6+(5-18*p._T(φ)+p._T(φ)*p._T(φ)+72*p._C(φ, sph)-58*sph.ei2())*
		math.Pow(A, 5)/120) + p.eastf
	north = p.scale*(p._M(φ, sph)-p._M(radian(p.latf), sph)+p._N(φ, sph)*math.Tan(φ)*
		(A*A/2+(5-p._T(φ)+9*p._C(φ, sph)+4*p._C(φ, sph)*p._C(φ, sph))*
			math.Pow(A, 4)/24+(61-58*p._T(φ)+p._T(φ)*p._T(φ)+600*
			p._C(φ, sph)-330*sph.ei2())*math.Pow(A, 6)/720)) + p.northf

	return east, north
}

func (transverseMercator) _M(φ float64, sph spheroid) float64 {
	return sph.A() * ((1-sph.e2()/4-3*sph.e4()/64-5*sph.e6()/256)*φ -
		(3*sph.e2()/8+3*sph.e4()/32+45*sph.e6()/1024)*math.Sin(2*φ) +
		(15*sph.e4()/256+45*sph.e6()/1024)*math.Sin(4*φ) -
		(35*sph.e6()/3072)*math.Sin(6*φ))
}

func (transverseMercator) _N(φ float64, sph spheroid) float64 {
	return sph.A() / math.Sqrt(1-sph.e2()*sin2(φ))
}

func (transverseMercator) _T(φ float64) float64 {
	return tan2(φ)
}

func (transverseMercator) _C(φ float64, sph spheroid) float64 {
	return sph.ei2() * cos2(φ)
}

type lambertConformalConic2SP struct {
	lonf, latf, lat1, lat2, eastf, northf float64
}

func (p lambertConformalConic2SP) ToLonLat(east, north float64, s Spheroid) (lon, lat float64) {
	sph := spheroid{a: s.A(), fi: s.Fi()}

	ρi := math.Sqrt(math.Pow(east-p.eastf, 2) + math.Pow(p._rho(radian(p.latf), sph)-(north-p.northf), 2))
	if p._n(sph) < 0 {
		ρi = -ρi
	}

	ti := math.Pow(ρi/(sph.A()*p._F(sph)), 1/p._n(sph))

	φ := math.Pi/2 - 2*math.Atan(ti)
	for i := 0; i < 5; i++ {
		φ = math.Pi/2 - 2*math.Atan(ti*math.Pow((1-sph.e()*math.Sin(φ))/(1+sph.e()*math.Sin(φ)), sph.e()/2))
	}

	λ := math.Atan((east-p.eastf)/(p._rho(radian(p.latf), sph)-(north-p.northf)))/p._n(sph) + radian(p.lonf)

	return degree(λ), degree(φ)
}

func (p lambertConformalConic2SP) FromLonLat(lon, lat float64, s Spheroid) (east, north float64) {
	sph := spheroid{a: s.A(), fi: s.Fi()}
	θ := p._n(sph) * (radian(lon) - radian(p.lonf))
	east = p.eastf + p._rho(radian(lat), sph)*math.Sin(θ)
	north = p.northf + p._rho(radian(p.latf), sph) - p._rho(radian(lat), sph)*math.Cos(θ)

	return east, north
}

func (p lambertConformalConic2SP) _t(φ float64, sph spheroid) float64 {
	return math.Tan(math.Pi/4-φ/2) /
		math.Pow((1-sph.e()*math.Sin(φ))/(1+sph.e()*math.Sin(φ)), sph.e()/2)
}

func (p lambertConformalConic2SP) _m(φ float64, sph spheroid) float64 {
	return math.Cos(φ) / math.Sqrt(1-sph.e2()*sin2(φ))
}

func (p lambertConformalConic2SP) _n(sph spheroid) float64 {
	if radian(p.lat1) == radian(p.lat2) {
		return math.Sin(radian(p.lat1))
	}

	return (math.Log(p._m(radian(p.lat1), sph)) - math.Log(p._m(radian(p.lat2), sph))) /
		(math.Log(p._t(radian(p.lat1), sph)) - math.Log(p._t(radian(p.lat2), sph)))
}

func (p lambertConformalConic2SP) _F(sph spheroid) float64 {
	return p._m(radian(p.lat1), sph) / (p._n(sph) * math.Pow(p._t(radian(p.lat1), sph), p._n(sph)))
}

func (p lambertConformalConic2SP) _rho(φ float64, sph spheroid) float64 {
	return sph.A() * p._F(sph) * math.Pow(p._t(φ, sph), p._n(sph))
}

type albersEqualAreaConic struct {
	lonf, latf, lat1, lat2, eastf, northf float64
}

func (p albersEqualAreaConic) ToLonLat(east, north float64, s Spheroid) (lon, lat float64) {
	sph := spheroid{a: s.A(), fi: s.Fi()}
	east -= p.eastf
	north -= p.northf
	ρi := math.Sqrt(east*east + math.Pow(p._rho(radian(p.latf), sph)-north, 2))
	qi := (p._C(sph) - ρi*ρi*p._n(sph)*p._n(sph)/sph.a2()) / p._n(sph)
	φ := math.Asin(qi / 2)

	for i := 0; i < 5; i++ {
		φ += math.Pow(1-sph.e2()*sin2(φ), 2) /
			(2 * math.Cos(φ)) * (qi/(1-sph.e2()) -
			math.Sin(φ)/(1-sph.e2()*sin2(φ)) +
			1/(2*sph.e())*math.Log((1-sph.e()*math.Sin(φ))/(1+sph.e()*math.Sin(φ))))
	}

	θ := math.Atan(east / (p._rho(radian(p.latf), sph) - north))

	return degree(radian(p.lonf) + θ/p._n(sph)), degree(φ)
}

func (p albersEqualAreaConic) FromLonLat(lon, lat float64, s Spheroid) (east, north float64) {
	sph := spheroid{a: s.A(), fi: s.Fi()}
	θ := p._n(sph) * (radian(lon) - radian(p.lonf))
	east = p.eastf + p._rho(radian(lat), sph)*math.Sin(θ)
	north = p.northf + p._rho(radian(p.latf), sph) - p._rho(radian(lat), sph)*math.Cos(θ)

	return east, north
}

func (p albersEqualAreaConic) _m(φ float64, sph spheroid) float64 {
	return math.Cos(φ) / math.Sqrt(1-sph.e2()*sin2(φ))
}

func (p albersEqualAreaConic) _q(φ float64, sph spheroid) float64 {
	return (1 - sph.e2()) * (math.Sin(φ)/(1-sph.e2()*sin2(φ)) -
		(1/(2*sph.e()))*math.Log((1-sph.e()*math.Sin(φ))/(1+sph.e()*math.Sin(φ))))
}

func (p albersEqualAreaConic) _n(sph spheroid) float64 {
	if radian(p.lat1) == radian(p.lat2) {
		return math.Sin(radian(p.lat1))
	}

	return (p._m(radian(p.lat1), sph)*p._m(radian(p.lat1), sph) -
		p._m(radian(p.lat2), sph)*p._m(radian(p.lat2), sph)) /
		(p._q(radian(p.lat2), sph) - p._q(radian(p.lat1), sph))
}

func (p albersEqualAreaConic) _C(sph spheroid) float64 {
	return p._m(radian(p.lat1), sph)*p._m(radian(p.lat1), sph) + p._n(sph)*p._q(radian(p.lat1), sph)
}

func (p albersEqualAreaConic) _rho(φ float64, sph spheroid) float64 {
	return sph.A() * math.Sqrt(p._C(sph)-p._n(sph)*p._q(φ, sph)) / p._n(sph)
}

type lambertAzimuthalEqualArea struct {
	latf, lonf, eastf, northf float64
}

func (p lambertAzimuthalEqualArea) ToLonLat(east, north float64, s Spheroid) (lon, lat float64) {
	sph := spheroid{a: s.A(), fi: s.Fi()}

	rho := math.Sqrt(math.Pow((east-p.eastf)/p._D(sph), 2) + math.Pow(p._D(sph)*(north-p.northf), 2))
	C := 2 * math.Asin(rho/(2*p._Rq(sph)))
	betaI := math.Asin((math.Cos(C) * math.Sin(p._beta0(sph))) +
		(p._D(sph) * (north - p.northf) * math.Sin(C) * math.Cos(p._beta0(sph)) / rho))

	if p.latf < 0 {
		betaI *= -1
	}

	rlat := betaI +
		((math.Pow(sph.e(), 2.0)/3.0 +
			31*math.Pow(sph.e(), 4.0)/180.0 +
			517*math.Pow(sph.e(), 6.0)/5040.0) *
			math.Sin(2*betaI)) +
		((23*math.Pow(sph.e(), 4.0)/360.0 +
			251*math.Pow(sph.e(), 6.0)/3780.0) *
			math.Sin(4*betaI)) +
		((761 * math.Pow(sph.e(), 6.0) / 45360.0) *
			math.Sin(6*betaI))

	lon = p.lonf +
		degree(math.Atan2((east-p.eastf)*math.Sin(C),
			(p._D(sph)*rho*
				math.Cos(p._beta0(sph))*
				math.Cos(C)-math.Pow(p._D(sph), 2)*(north-p.northf)*math.Sin(p._beta0(sph))*math.Sin(C))))

	return lon, degree(rlat)
}

func (p lambertAzimuthalEqualArea) FromLonLat(lon, lat float64, s Spheroid) (east, north float64) {
	sph := spheroid{a: s.A(), fi: s.Fi()}

	beta := math.Asin(p._q(lat, sph) / p._qp(sph))
	B := p._Rq(sph) * math.Sqrt(2/(1+math.Sin(p._beta0(sph))*math.Sin(beta)+
		math.Cos(p._beta0(sph))*math.Cos(beta)*math.Cos(radian(lon-p.lonf))))

	return (p.eastf + B*p._D(sph)*math.Cos(beta)*math.Sin(radian(lon-p.lonf))),
		(p.northf + B/p._D(sph)*(math.Cos(p._beta0(sph))*math.Sin(beta)-
			math.Sin(p._beta0(sph))*math.Cos(beta)*math.Cos(radian(lon-p.lonf))))
}

func (p lambertAzimuthalEqualArea) _q(lat float64, sph spheroid) float64 {
	return (1 - sph.e2()) *
		((math.Sin(radian(lat)) / (1 - sph.e2()*sin2(radian(lat)))) -
			((1 / (2 * sph.e())) * math.Log((1-sph.e()*math.Sin(radian(lat)))/(1+sph.e()*math.Sin(radian(lat))))))
}

func (p lambertAzimuthalEqualArea) _qp(sph spheroid) float64 {
	return (1 - sph.e2()) * ((1 / (1 - sph.e2())) - ((1 / (2 * sph.e())) * math.Log((1-sph.e())/(1+sph.e()))))
}

func (p lambertAzimuthalEqualArea) _q0(sph spheroid) float64 {
	return (1 - sph.e2()) * ((math.Sin(radian(p.latf)) / (1 - sph.e2()*sin2(radian(p.latf)))) -
		((1 / (2 * sph.e())) * math.Log((1-sph.e()*math.Sin(radian(p.latf)))/(1+sph.e()*math.Sin(radian(p.latf))))))
}

func (p lambertAzimuthalEqualArea) _beta0(sph spheroid) float64 {
	return math.Asin(p._q0(sph) / p._qp(sph))
}

func (p lambertAzimuthalEqualArea) _Rq(sph spheroid) float64 {
	return sph.A() * math.Pow(p._qp(sph)/2, 0.5)
}

func (p lambertAzimuthalEqualArea) _D(sph spheroid) float64 {
	return sph.A() * (math.Cos(radian(p.latf)) / math.Sqrt(1-sph.e2()*math.Pow(math.Sin(radian(p.latf)), 2))) /
		(p._Rq(sph) * math.Cos(p._beta0(sph)))
}
