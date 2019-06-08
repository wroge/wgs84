package wgs84

import (
	"fmt"
	"strconv"

	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
	"github.com/wroge/wgs84/transformation"
)

type CoordinateReferenceSystem struct {
	Spheroid       spheroid.Spheroid
	Transformation transformation.Transformation
	System         system.System
}

func XYZ() CoordinateReferenceSystem {
	return System(system.XYZ())
}

func LonLat() CoordinateReferenceSystem {
	return System(system.LonLat())
}

func WebMercator() CoordinateReferenceSystem {
	return System(system.WebMercator())
}

func UTM(zone float64, northern bool) CoordinateReferenceSystem {
	return System(system.UTM(zone, northern))
}

func System(sys system.System) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       spheroid.WGS84(),
		Transformation: transformation.WGS84(),
		System:         sys,
	}
}

type Func func(a, b, c float64) (a2, b2, c2 float64)

func (f Func) Round(places uint) Func {
	r := func(v float64) float64 {
		v, _ = strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(int(places))+"f", v), 64)
		if v == -0 {
			return 0
		}
		return v
	}
	return func(a, b, c float64) (x, y, z float64) {
		a, b, c = f(a, b, c)
		return r(a), r(b), r(c)
	}
}

func (crs CoordinateReferenceSystem) ToSystem(sys system.System) Func {
	return func(a, b, c float64) (a2, b2, c2 float64) {
		x, y, z := crs.System.ToXYZ(a, b, c, crs.Spheroid)
		return sys.FromXYZ(x, y, z, crs.Spheroid)
	}
}

func (crs CoordinateReferenceSystem) FromSystem(sys system.System) Func {
	return func(a, b, c float64) (a2, b2, c2 float64) {
		x, y, z := sys.ToXYZ(a, b, c, crs.Spheroid)
		return crs.System.FromXYZ(x, y, z, crs.Spheroid)
	}
}

func (crs CoordinateReferenceSystem) ToXYZ() Func {
	return crs.ToSystem(system.XYZ())
}

func (crs CoordinateReferenceSystem) ToLonLat() Func {
	return crs.ToSystem(system.LonLat())
}

func (crs CoordinateReferenceSystem) ToWebMercator() Func {
	return crs.ToSystem(system.WebMercator())
}

func (crs CoordinateReferenceSystem) ToTransverseMercator(lonf, latf, scale, eastf, northf float64) Func {
	return crs.ToSystem(system.TransverseMercator(lonf, latf, scale, eastf, northf))
}

func (crs CoordinateReferenceSystem) ToUTM(zone float64, northern bool) Func {
	return crs.ToSystem(system.UTM(zone, northern))
}

func (crs CoordinateReferenceSystem) ToGK(zone float64) Func {
	return crs.ToSystem(system.GK(zone))
}

func (crs CoordinateReferenceSystem) ToMercator(lonf, scale, eastf, northf float64) Func {
	return crs.ToSystem(system.Mercator(lonf, scale, eastf, northf))
}

func (crs CoordinateReferenceSystem) ToLambertConformalConic1SP(lonf, latf, scale, eastf, northf float64) Func {
	return crs.ToSystem(system.LambertConformalConic1SP(lonf, latf, scale, eastf, northf))
}

func (crs CoordinateReferenceSystem) ToLambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.ToSystem(system.LambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf))
}

func (crs CoordinateReferenceSystem) ToEquidistantConic(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.ToSystem(system.EquidistantConic(lonf, latf, lat1, lat2, eastf, northf))
}

func (crs CoordinateReferenceSystem) ToAlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.ToSystem(system.AlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf))
}

func (crs CoordinateReferenceSystem) FromXYZ() Func {
	return crs.FromSystem(system.XYZ())
}

func (crs CoordinateReferenceSystem) FromLonLat() Func {
	return crs.FromSystem(system.LonLat())
}

func (crs CoordinateReferenceSystem) FromWebMercator() Func {
	return crs.FromSystem(system.WebMercator())
}

func (crs CoordinateReferenceSystem) FromTransverseMercator(lonf, latf, scale, eastf, northf float64) Func {
	return crs.FromSystem(system.TransverseMercator(lonf, latf, scale, eastf, northf))
}

func (crs CoordinateReferenceSystem) FromUTM(zone float64, northern bool) Func {
	return crs.FromSystem(system.UTM(zone, northern))
}

func (crs CoordinateReferenceSystem) FromGK(zone float64) Func {
	return crs.FromSystem(system.GK(zone))
}

func (crs CoordinateReferenceSystem) FromMercator(lonf, scale, eastf, northf float64) Func {
	return crs.FromSystem(system.Mercator(lonf, scale, eastf, northf))
}

func (crs CoordinateReferenceSystem) FromLambertConformalConic1SP(lonf, latf, scale, eastf, northf float64) Func {
	return crs.FromSystem(system.LambertConformalConic1SP(lonf, latf, scale, eastf, northf))
}

func (crs CoordinateReferenceSystem) FromLambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.FromSystem(system.LambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf))
}

func (crs CoordinateReferenceSystem) FromEquidistantConic(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.FromSystem(system.EquidistantConic(lonf, latf, lat1, lat2, eastf, northf))
}

func (crs CoordinateReferenceSystem) FromAlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.FromSystem(system.AlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf))
}

func (crs CoordinateReferenceSystem) To(to CoordinateReferenceSystem) Func {
	return func(a, b, c float64) (a2, b2, c2 float64) {
		x, y, z := crs.System.ToXYZ(a, b, c, crs.Spheroid)
		x0, y0, z0 := crs.Transformation.ToWGS84(x, y, z)
		x2, y2, z2 := to.Transformation.FromWGS84(x0, y0, z0)
		return to.System.FromXYZ(x2, y2, z2, to.Spheroid)
	}
}

func (crs CoordinateReferenceSystem) From(from CoordinateReferenceSystem) Func {
	return from.To(crs)
}

func ToXYZ() Func {
	return CoordinateReferenceSystem{}.ToXYZ()
}

func ToLonLat() Func {
	return CoordinateReferenceSystem{}.ToLonLat()
}

func ToWebMercator() Func {
	return CoordinateReferenceSystem{}.ToWebMercator()
}

func ToUTM(zone float64, northern bool) Func {
	return CoordinateReferenceSystem{}.ToUTM(zone, northern)
}

func FromXYZ() Func {
	return CoordinateReferenceSystem{}.FromXYZ()
}

func FromLonLat() Func {
	return CoordinateReferenceSystem{}.FromLonLat()
}

func FromWebMercator() Func {
	return CoordinateReferenceSystem{}.FromWebMercator()
}

func FromUTM(zone float64, northern bool) Func {
	return CoordinateReferenceSystem{}.FromUTM(zone, northern)
}

func To(to CoordinateReferenceSystem) Func {
	return CoordinateReferenceSystem{}.To(to)
}

func From(from CoordinateReferenceSystem) Func {
	return CoordinateReferenceSystem{}.From(from)
}
