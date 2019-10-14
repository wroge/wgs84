package wgs84

import (
	"math"
)

// XYZ is a geocentric CoordinateReferenceSystem which is by default
// similar to WGS84 and the EPSG-Code 4978.
type XYZ struct {
	GeodeticDatum GeodeticDatum
	Area          Area
}

func (crs XYZ) Contains(lon, lat float64) bool {
	if crs.Area != nil && !crs.Area.Contains(lon, lat) {
		return false
	}
	if crs.GeodeticDatum != nil && !crs.GeodeticDatum.Contains(lon, lat) {
		return false
	}
	if math.Abs(lon) > 180 || math.Abs(lat) > 90 {
		return false
	}
	return true
}

func (crs XYZ) To(to CoordinateReferenceSystem) Func {
	return Transform(crs, to)
}

func (crs XYZ) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs XYZ) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs XYZ) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs XYZ) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs XYZ) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	return a, b, c
}

func (crs XYZ) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	return x, y, z
}

// LonLat is a geographic CoordinateReferenceSystem which is by default
// similar to WGS84 and the EPSG-Code 4326.
type LonLat struct {
	GeodeticDatum GeodeticDatum
	Area          Area
}

func (crs LonLat) Contains(lon, lat float64) bool {
	if crs.Area != nil && !crs.Area.Contains(lon, lat) {
		return false
	}
	if crs.GeodeticDatum != nil && !crs.GeodeticDatum.Contains(lon, lat) {
		return false
	}
	if math.Abs(lon) > 180 || math.Abs(lat) > 90 {
		return false
	}
	return true
}

func (crs LonLat) To(to CoordinateReferenceSystem) Func {
	return Transform(crs, to)
}

func (crs LonLat) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs LonLat) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs LonLat) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs LonLat) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs LonLat) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	x = (crs._N(radian(b), s) + c) * math.Cos(radian(a)) * math.Cos(radian(b))
	y = (crs._N(radian(b), s) + c) * math.Cos(radian(b)) * math.Sin(radian(a))
	z = (crs._N(radian(b), s)*math.Pow(s.MajorAxis()*(1-s.F()), 2)/(s.A2()) + c) * math.Sin(radian(b))
	return x, y, z
}

func (crs LonLat) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	sd := math.Sqrt(x*x + y*y)
	T := math.Atan(z * s.MajorAxis() / (sd * s.B()))
	B := math.Atan((z + s.E2()*(s.A2())/s.B()*
		math.Pow(math.Sin(T), 3)) / (sd - s.E2()*s.MajorAxis()*math.Pow(math.Cos(T), 3)))
	c = sd/math.Cos(B) - crs._N(B, s)
	a = degree(math.Atan2(y, x))
	b = degree(B)
	return a, b, c
}

func (crs LonLat) _N(φ float64, s Spheroid) float64 {
	return s.MajorAxis() / math.Sqrt(1-s.E2()*math.Pow(math.Sin(φ), 2))
}

// Projection is a wrapper for projected CoordinateReferenceSystems.
type Projection struct {
	GeodeticDatum        GeodeticDatum
	CoordinateProjection CoordinateProjection
	Area                 Area
}

func (crs Projection) Contains(lon, lat float64) bool {
	if crs.Area != nil && !crs.Area.Contains(lon, lat) {
		return false
	}
	if crs.GeodeticDatum != nil && !crs.GeodeticDatum.Contains(lon, lat) {
		return false
	}
	if math.Abs(lon) > 180 || math.Abs(lat) > 90 {
		return false
	}
	return true
}

func (crs Projection) To(to CoordinateReferenceSystem) Func {
	return Transform(crs, to)
}

func (crs Projection) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs Projection) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs Projection) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs Projection) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs Projection) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	if crs.CoordinateProjection == nil {
		return WebMercator{
			GeodeticDatum: crs.GeodeticDatum,
		}.ToXYZ(a, b, c, gs)
	}
	s := spheroid(gs, crs.GeodeticDatum)
	a, b = crs.CoordinateProjection.ToLonLat(a, b, s)
	return LonLat{
		GeodeticDatum: crs.GeodeticDatum,
	}.ToXYZ(a, b, c, s)
}

func (crs Projection) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	if crs.CoordinateProjection == nil {
		return WebMercator{
			GeodeticDatum: crs.GeodeticDatum,
		}.FromXYZ(x, y, z, gs)
	}
	s := spheroid(gs, crs.GeodeticDatum)
	a, b, c = LonLat{
		GeodeticDatum: crs.GeodeticDatum,
	}.FromXYZ(x, y, z, s)
	a, b = crs.CoordinateProjection.FromLonLat(a, b, s)
	return a, b, c
}

// WebMercator is a projected CoordinateReferenceSystem which is by default
// similar to WGS84 and the EPSG-Code 3857.
type WebMercator struct {
	GeodeticDatum GeodeticDatum
	Area          Area
}

func (crs WebMercator) Contains(lon, lat float64) bool {
	if crs.Area != nil && !crs.Area.Contains(lon, lat) {
		return false
	}
	if crs.GeodeticDatum != nil && !crs.GeodeticDatum.Contains(lon, lat) {
		return false
	}
	if math.Abs(lon) > 180 || math.Abs(lat) > 85.06 {
		return false
	}
	return true
}

func (crs WebMercator) To(to CoordinateReferenceSystem) Func {
	return Transform(crs, to)
}

func (crs WebMercator) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs WebMercator) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs WebMercator) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs WebMercator) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs WebMercator) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs WebMercator) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs WebMercator) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	lon = degree(east / s.MajorAxis())
	lat = math.Atan(math.Exp(north/s.MajorAxis()))*degree(1)*2 - 90
	return lon, lat
}

func (crs WebMercator) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	east = radian(lon) * s.MajorAxis()
	north = math.Log(math.Tan(radian((90+lat)/2))) * s.MajorAxis()
	return east, north
}

// UTM is a projected CoordinateReferenceSystem. It is a TransverseMercator
// System with 6 degrees zone width, 0.9996 Scale on the central meridian
// and 10000000 false northing on the southern hemisphere. It is similar to
// the EPSG-Codes 32601 - 32660 on the northern hemisphere and 32701 - 32760
// on the southern hemisphere.
func UTM(zone float64, northern bool) TransverseMercator {
	northf := 0.0
	if !northern {
		northf = 10000000
	}
	return TransverseMercator{
		GeodeticDatum: Datum{
			GeodeticSpheroid: Spheroid{},
		},
		Lonf:   zone*6 - 183,
		Latf:   0,
		Scale:  0.9996,
		Eastf:  500000,
		Northf: northf,
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < zone*6-186 || lon > zone*6-180 {
				return false
			}
			if northern && (lat < 0 || lat > 84) {
				return false
			}
			if !northern && (lat > 0 || lat < -80) {
				return false
			}
			return true
		}),
	}
}

// TransverseMercator is a projected CoordinateReferenceSystem. By default
// the GeodeticDatum is similar to WGS84.
type TransverseMercator struct {
	Lonf, Latf, Scale, Eastf, Northf float64
	GeodeticDatum                    GeodeticDatum
	Area                             Area
}

func (crs TransverseMercator) Contains(lon, lat float64) bool {
	if crs.Area != nil && !crs.Area.Contains(lon, lat) {
		return false
	}
	if crs.GeodeticDatum != nil && !crs.GeodeticDatum.Contains(lon, lat) {
		return false
	}
	if math.Abs(lon) > 180 || math.Abs(lat) > 90 {
		return false
	}
	return true
}

func (crs TransverseMercator) To(to CoordinateReferenceSystem) Func {
	return Transform(crs, to)
}

func (crs TransverseMercator) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs TransverseMercator) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs TransverseMercator) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs TransverseMercator) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs TransverseMercator) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs TransverseMercator) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs TransverseMercator) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	east -= crs.Eastf
	north -= crs.Northf
	Mi := crs._M(radian(crs.Latf), s) + north/crs.Scale
	μ := Mi / (s.MajorAxis() * (1 - s.E2()/4 - 3*s.E4()/64 - 5*s.E6()/256))
	φ1 := μ + (3*s.Ei()/2-27*s.Ei3()/32)*math.Sin(2*μ) +
		(21*s.Ei2()/16-55*s.Ei4()/32)*math.Sin(4*μ) +
		(151*s.Ei3()/96)*math.Sin(6*μ) +
		(1097*s.Ei4()/512)*math.Sin(8*μ)
	R1 := s.MajorAxis() * (1 - s.E2()) / math.Pow(1-s.E2()*sin2(φ1), 3/2)
	D := east / (crs._N(φ1, s) * crs.Scale)
	φ := φ1 - (crs._N(φ1, s)*math.Tan(φ1)/R1)*(D*D/2-(5+3*crs._T(φ1)+10*
		crs._C(φ1, s)-4*crs._C(φ1, s)*crs._C(φ1, s)-9*s.Ei2())*
		math.Pow(D, 4)/24+(61+90*crs._T(φ1)+298*crs._C(φ1, s)+45*crs._T(φ1)*
		crs._T(φ1)-252*s.Ei2()-3*crs._C(φ1, s)*crs._C(φ1, s))*
		math.Pow(D, 6)/720)
	λ := radian(crs.Lonf) + (D-(1+2*crs._T(φ1)+crs._C(φ1, s))*D*D*D/6+(5-2*crs._C(φ1, s)+
		28*crs._T(φ1)-3*crs._C(φ1, s)*crs._C(φ1, s)+8*s.Ei2()+24*crs._T(φ1)*crs._T(φ1))*
		math.Pow(D, 5)/120)/math.Cos(φ1)
	return degree(λ), degree(φ)
}

func (crs TransverseMercator) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	φ := radian(lat)
	A := (radian(lon) - radian(crs.Lonf)) * math.Cos(φ)
	east = crs.Scale*crs._N(φ, s)*(A+(1-crs._T(φ)+crs._C(φ, s))*
		math.Pow(A, 3)/6+(5-18*crs._T(φ)+crs._T(φ)*crs._T(φ)+72*crs._C(φ, s)-58*s.Ei2())*
		math.Pow(A, 5)/120) + crs.Eastf
	north = crs.Scale*(crs._M(φ, s)-crs._M(radian(crs.Latf), s)+crs._N(φ, s)*math.Tan(φ)*
		(A*A/2+(5-crs._T(φ)+9*crs._C(φ, s)+4*crs._C(φ, s)*crs._C(φ, s))*
			math.Pow(A, 4)/24+(61-58*crs._T(φ)+crs._T(φ)*crs._T(φ)+600*
			crs._C(φ, s)-330*s.Ei2())*math.Pow(A, 6)/720)) + crs.Northf
	return east, north
}

func (TransverseMercator) _M(φ float64, s Spheroid) float64 {
	return s.MajorAxis() * ((1-s.E2()/4-3*s.E4()/64-5*s.E6()/256)*φ -
		(3*s.E2()/8+3*s.E4()/32+45*s.E6()/1024)*math.Sin(2*φ) +
		(15*s.E4()/256+45*s.E6()/1024)*math.Sin(4*φ) -
		(35*s.E6()/3072)*math.Sin(6*φ))
}

func (TransverseMercator) _N(φ float64, s Spheroid) float64 {
	return s.MajorAxis() / math.Sqrt(1-s.E2()*sin2(φ))
}

func (TransverseMercator) _T(φ float64) float64 {
	return tan2(φ)
}

func (TransverseMercator) _C(φ float64, s Spheroid) float64 {
	return s.Ei2() * cos2(φ)
}

// LambertConformalConic2SP is a projected CoordinateReferenceSystem. By
// default the GeodeticDatum is similar to WGS84.
type LambertConformalConic2SP struct {
	Lonf, Latf, Lat1, Lat2, Eastf, Northf float64
	GeodeticDatum                         GeodeticDatum
	Area                                  Area
}

func (crs LambertConformalConic2SP) Contains(lon, lat float64) bool {
	if crs.Area != nil && !crs.Area.Contains(lon, lat) {
		return false
	}
	if crs.GeodeticDatum != nil && !crs.GeodeticDatum.Contains(lon, lat) {
		return false
	}
	if math.Abs(lon) > 180 || math.Abs(lat) > 90 {
		return false
	}
	return true
}

func (crs LambertConformalConic2SP) To(to CoordinateReferenceSystem) Func {
	return Transform(crs, to)
}

func (crs LambertConformalConic2SP) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs LambertConformalConic2SP) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs LambertConformalConic2SP) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs LambertConformalConic2SP) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs LambertConformalConic2SP) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs LambertConformalConic2SP) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs LambertConformalConic2SP) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	ρi := math.Sqrt(math.Pow(east-crs.Eastf, 2) + math.Pow(crs._ρ(radian(crs.Latf), s)-(north-crs.Northf), 2))
	if crs._n(s) < 0 {
		ρi = -ρi
	}
	ti := math.Pow(ρi/(s.MajorAxis()*crs._F(s)), 1/crs._n(s))
	φ := math.Pi/2 - 2*math.Atan(ti)
	for i := 0; i < 5; i++ {
		φ = math.Pi/2 - 2*math.Atan(ti*math.Pow((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ)), s.E()/2))
	}
	λ := math.Atan((east-crs.Eastf)/(crs._ρ(radian(crs.Latf), s)-(north-crs.Northf)))/crs._n(s) + radian(crs.Lonf)
	return degree(λ), degree(φ)
}

func (crs LambertConformalConic2SP) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	θ := crs._n(s) * (radian(lon) - radian(crs.Lonf))
	east = crs.Eastf + crs._ρ(radian(lat), s)*math.Sin(θ)
	north = crs.Northf + crs._ρ(radian(crs.Latf), s) - crs._ρ(radian(lat), s)*math.Cos(θ)
	return east, north
}

func (crs LambertConformalConic2SP) _t(φ float64, s Spheroid) float64 {
	return math.Tan(math.Pi/4-φ/2) /
		math.Pow((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ)), s.E()/2)
}

func (crs LambertConformalConic2SP) _m(φ float64, s Spheroid) float64 {
	return math.Cos(φ) / math.Sqrt(1-s.E2()*sin2(φ))
}

func (crs LambertConformalConic2SP) _n(s Spheroid) float64 {
	if radian(crs.Lat1) == radian(crs.Lat2) {
		return math.Sin(radian(crs.Lat1))
	}
	return (math.Log(crs._m(radian(crs.Lat1), s)) - math.Log(crs._m(radian(crs.Lat2), s))) /
		(math.Log(crs._t(radian(crs.Lat1), s)) - math.Log(crs._t(radian(crs.Lat2), s)))
}

func (crs LambertConformalConic2SP) _F(s Spheroid) float64 {
	return crs._m(radian(crs.Lat1), s) / (crs._n(s) * math.Pow(crs._t(radian(crs.Lat1), s), crs._n(s)))
}

func (crs LambertConformalConic2SP) _ρ(φ float64, s Spheroid) float64 {
	return s.MajorAxis() * crs._F(s) * math.Pow(crs._t(φ, s), crs._n(s))
}

// AlbersEqualAreaConic is a projected CoordinateReferenceSystem. By
// default the GeodeticDatum is similar to WGS84.
type AlbersEqualAreaConic struct {
	Lonf, Latf, Lat1, Lat2, Eastf, Northf float64
	GeodeticDatum                         GeodeticDatum
	Area                                  Area
}

func (crs AlbersEqualAreaConic) Contains(lon, lat float64) bool {
	if crs.Area != nil && !crs.Area.Contains(lon, lat) {
		return false
	}
	if crs.GeodeticDatum != nil && !crs.GeodeticDatum.Contains(lon, lat) {
		return false
	}
	if math.Abs(lon) > 180 || math.Abs(lat) > 90 {
		return false
	}
	return true
}

func (crs AlbersEqualAreaConic) To(to CoordinateReferenceSystem) Func {
	return Transform(crs, to)
}

func (crs AlbersEqualAreaConic) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs AlbersEqualAreaConic) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs AlbersEqualAreaConic) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs AlbersEqualAreaConic) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs AlbersEqualAreaConic) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs AlbersEqualAreaConic) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs AlbersEqualAreaConic) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	east -= crs.Eastf
	north -= crs.Northf
	ρi := math.Sqrt(east*east + math.Pow(crs._ρ(radian(crs.Latf), s)-north, 2))
	qi := (crs._C(s) - ρi*ρi*crs._n(s)*crs._n(s)/s.A2()) / crs._n(s)
	φ := math.Asin(qi / 2)
	for i := 0; i < 5; i++ {
		φ += math.Pow(1-s.E2()*sin2(φ), 2) /
			(2 * math.Cos(φ)) * (qi/(1-s.E2()) -
			math.Sin(φ)/(1-s.E2()*sin2(φ)) +
			1/(2*s.E())*math.Log((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ))))
	}
	θ := math.Atan(east / (crs._ρ(radian(crs.Latf), s) - north))
	return degree(radian(crs.Lonf) + θ/crs._n(s)), degree(φ)
}

func (crs AlbersEqualAreaConic) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	θ := crs._n(s) * (radian(lon) - radian(crs.Lonf))
	east = crs.Eastf + crs._ρ(radian(lat), s)*math.Sin(θ)
	north = crs.Northf + crs._ρ(radian(crs.Latf), s) - crs._ρ(radian(lat), s)*math.Cos(θ)
	return east, north
}

func (crs AlbersEqualAreaConic) _m(φ float64, s Spheroid) float64 {
	return math.Cos(φ) / math.Sqrt(1-s.E2()*sin2(φ))
}

func (crs AlbersEqualAreaConic) _q(φ float64, s Spheroid) float64 {
	return (1 - s.E2()) * (math.Sin(φ)/(1-s.E2()*sin2(φ)) -
		(1/(2*s.E()))*math.Log((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ))))
}

func (crs AlbersEqualAreaConic) _n(s Spheroid) float64 {
	if radian(crs.Lat1) == radian(crs.Lat2) {
		return math.Sin(radian(crs.Lat1))
	}
	return (crs._m(radian(crs.Lat1), s)*crs._m(radian(crs.Lat1), s) -
		crs._m(radian(crs.Lat2), s)*crs._m(radian(crs.Lat2), s)) /
		(crs._q(radian(crs.Lat2), s) - crs._q(radian(crs.Lat1), s))
}

func (crs AlbersEqualAreaConic) _C(s Spheroid) float64 {
	return crs._m(radian(crs.Lat1), s)*crs._m(radian(crs.Lat1), s) + crs._n(s)*crs._q(radian(crs.Lat1), s)
}

func (crs AlbersEqualAreaConic) _ρ(φ float64, s Spheroid) float64 {
	return s.MajorAxis() * math.Sqrt(crs._C(s)-crs._n(s)*crs._q(φ, s)) / crs._n(s)
}
