// Package wgs84 provides multiple coordinate
// conversions and transformations.
package wgs84

import (
	"fmt"
	"strconv"

	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
	"github.com/wroge/wgs84/transformation"
)

// The Transformation interface is implemented by
// github.com/wroge/wgs84/transformation transformation.Transformation.
type Transformation interface {
	ToWGS84(x, y, z float64) (x0, y0, z0 float64)
	FromWGS84(x0, y0, z0 float64) (x, y, z float64)
}

// Spheroid is an unnamed interface type implemented by
// github.com/wroge/wgs84/spheroid spheroid.Spheroid.
type Spheroid = interface {
	A() float64
	Fi() float64
}

// The System interface is implemented by
// github.com/wroge/wgs84/system System.
type System interface {
	ToXYZ(a, b, c float64, sph Spheroid) (x, y, z float64)
	FromXYZ(x, y, z float64, sph Spheroid) (a, b, c float64)
}

// A CoordinateReferenceSystem contains a geodetic
// datum and a coordinate system (System).
// A geodetic datum is composed of the two components
// Spheroid and Transformation. A Transformation is a
// struct for transforming geocentric coordinates from
// and to the WGS84 system.
type CoordinateReferenceSystem struct {
	Spheroid       Spheroid
	Transformation Transformation
	System         System
}

// XYZ is the geocentric WGS84 system. It is equivalent
// to the EPSG-Code: 4978. The unit is meter.
func XYZ() CoordinateReferenceSystem {
	return WithSystem(system.XYZ())
}

// LonLat is the geographic WGS84 system. It is equivalent
// to the EPSG-Code: 4326. The unit is degrees.
func LonLat() CoordinateReferenceSystem {
	return WithSystem(system.LonLat())
}

// WebMercator is a projected WGS84 system often used by
// web maps like google maps. It is equivalent to the
// EPSG-Code: 3587. The unit is meter.
func WebMercator() CoordinateReferenceSystem {
	return WithSystem(system.WebMercator())
}

// UTM is a transverse mercator projection in the WGS84 system.
// It has 6 degrees zone width and a scale factor of 0.9996 at
// the central meridian. The unit is meter.
func UTM(zone float64, northern bool) CoordinateReferenceSystem {
	return WithSystem(system.UTM(zone, northern))
}

// WithSystem it is possible to use any System in the
// WGS84 geodetic datum. Only the most commonly used System's
// for WGS84 are implemented directly in this package.
func WithSystem(sys System) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       spheroid.WGS84(),
		Transformation: transformation.WGS84(),
		System:         sys,
	}
}

// The Func type is the result of all conversions and transformations.
// With this type it is possible to spice up the result, for
// example by rounding functions or changing the order of the
// parameters.
type Func func(a, b, c float64) (a2, b2, c2 float64)

// Round the result of a Func to a specified
// decimal place. It is not a rounding to a certain accuracy,
// as some coordinates are defined as degrees (LonLat) and others
// as meters (projected and geocentric).
func (f Func) Round(places uint) Func {
	if f == nil {
		return nil
	}
	r := func(v float64) float64 {
		v, _ = strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(int(places))+"f", v), 64)
		if v == -0 {
			return 0
		}
		return v
	}
	return func(a, b, c float64) (a2, b2, c2 float64) {
		a, b, c = f(a, b, c)
		return r(a), r(b), r(c)
	}
}

// SwitchIn changes the order of the first two incoming
// parameters. By default geographic coordinates where
// used in the order longitude, latitude, height. With
// this function it is possible to insert the parameters
// in the order latitude, longitude, height.
func (f Func) SwitchIn() Func {
	if f == nil {
		return nil
	}
	return func(a, b, c float64) (a2, b2, c2 float64) {
		return f(b, a, c)
	}
}

// SwitchOut changes the order of the first two resulting
// parameters. By default geographic coordinates where
// used in the order longitude, latitude, height. With
// this function it is possible to get the result
// in the order latitude, longitude, height.
func (f Func) SwitchOut() Func {
	if f == nil {
		return nil
	}
	return func(a, b, c float64) (a2, b2, c2 float64) {
		a, b, c = f(a, b, c)
		return b, a, c
	}
}

// ToSystem provides the conversion of a CoordinateReferenceSystem
// to a coordinate system implementing the System interface.
func (crs CoordinateReferenceSystem) ToSystem(sys System) Func {
	if crs.System == nil || crs.Spheroid == nil || crs.Transformation == nil {
		return nil
	}
	return func(a, b, c float64) (a2, b2, c2 float64) {
		x, y, z := crs.System.ToXYZ(a, b, c, crs.Spheroid)
		return sys.FromXYZ(x, y, z, crs.Spheroid)
	}
}

// FromSystem provides the conversion from a coordinate system
// implementing the System interface to a CoordinateReferenceSystem.
func (crs CoordinateReferenceSystem) FromSystem(sys System) Func {
	if crs.System == nil || crs.Spheroid == nil || crs.Transformation == nil {
		return nil
	}
	return func(a, b, c float64) (a2, b2, c2 float64) {
		x, y, z := sys.ToXYZ(a, b, c, crs.Spheroid)
		return crs.System.FromXYZ(x, y, z, crs.Spheroid)
	}
}

// ToXYZ is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are x, y and z in units of meters.
func (crs CoordinateReferenceSystem) ToXYZ() Func {
	return crs.ToSystem(system.XYZ())
}

// ToLonLat is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are lon, lat, h in units of degrees.
func (crs CoordinateReferenceSystem) ToLonLat() Func {
	return crs.ToSystem(system.LonLat())
}

// ToWebMercator is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) ToWebMercator() Func {
	return crs.ToSystem(system.WebMercator())
}

// ToTransverseMercator is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) ToTransverseMercator(lonf, latf, scale, eastf, northf float64) Func {
	return crs.ToSystem(system.TransverseMercator(lonf, latf, scale, eastf, northf))
}

// ToUTM is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) ToUTM(zone float64, northern bool) Func {
	return crs.ToSystem(system.UTM(zone, northern))
}

// ToGK is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) ToGK(zone float64) Func {
	return crs.ToSystem(system.GK(zone))
}

// ToMercator is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) ToMercator(lonf, scale, eastf, northf float64) Func {
	return crs.ToSystem(system.Mercator(lonf, scale, eastf, northf))
}

// ToLambertConformalConic1SP is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) ToLambertConformalConic1SP(lonf, latf, scale, eastf, northf float64) Func {
	return crs.ToSystem(system.LambertConformalConic1SP(lonf, latf, scale, eastf, northf))
}

// ToLambertConformalConic2SP is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) ToLambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.ToSystem(system.LambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf))
}

// ToEquidistantConic is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) ToEquidistantConic(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.ToSystem(system.EquidistantConic(lonf, latf, lat1, lat2, eastf, northf))
}

// ToAlbersEqualAreaConic is a coordinate conversion within the WGS84 geodetic datum.
// The resulting parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) ToAlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.ToSystem(system.AlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf))
}

// FromXYZ is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are x, y and z in units of meters.
func (crs CoordinateReferenceSystem) FromXYZ() Func {
	return crs.FromSystem(system.XYZ())
}

// FromLonLat is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are lon, lat, h in units of degrees.
func (crs CoordinateReferenceSystem) FromLonLat() Func {
	return crs.FromSystem(system.LonLat())
}

// FromWebMercator is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) FromWebMercator() Func {
	return crs.FromSystem(system.WebMercator())
}

// FromTransverseMercator is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) FromTransverseMercator(lonf, latf, scale, eastf, northf float64) Func {
	return crs.FromSystem(system.TransverseMercator(lonf, latf, scale, eastf, northf))
}

// FromUTM is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) FromUTM(zone float64, northern bool) Func {
	return crs.FromSystem(system.UTM(zone, northern))
}

// FromGK is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) FromGK(zone float64) Func {
	return crs.FromSystem(system.GK(zone))
}

// FromMercator is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) FromMercator(lonf, scale, eastf, northf float64) Func {
	return crs.FromSystem(system.Mercator(lonf, scale, eastf, northf))
}

// FromLambertConformalConic1SP is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) FromLambertConformalConic1SP(lonf, latf, scale, eastf, northf float64) Func {
	return crs.FromSystem(system.LambertConformalConic1SP(lonf, latf, scale, eastf, northf))
}

// FromLambertConformalConic2SP is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) FromLambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.FromSystem(system.LambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf))
}

// FromEquidistantConic is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) FromEquidistantConic(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.FromSystem(system.EquidistantConic(lonf, latf, lat1, lat2, eastf, northf))
}

// FromAlbersEqualAreaConic is a coordinate conversion within the WGS84 geodetic datum.
// The incoming parameters are east, north and h in units of meters.
func (crs CoordinateReferenceSystem) FromAlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf float64) Func {
	return crs.FromSystem(system.AlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf))
}

// To transforms from one CoordinateReferenceSystem to another.
func (crs CoordinateReferenceSystem) To(to CoordinateReferenceSystem) Func {
	if crs.System == nil || crs.Spheroid == nil || crs.Transformation == nil {
		return nil
	}
	if to.System == nil || to.Spheroid == nil || to.Transformation == nil {
		return nil
	}
	return func(a, b, c float64) (a2, b2, c2 float64) {
		x, y, z := crs.System.ToXYZ(a, b, c, crs.Spheroid)
		x0, y0, z0 := crs.Transformation.ToWGS84(x, y, z)
		x2, y2, z2 := to.Transformation.FromWGS84(x0, y0, z0)
		return to.System.FromXYZ(x2, y2, z2, to.Spheroid)
	}
}

// From is the reverse of the To func.
func (crs CoordinateReferenceSystem) From(from CoordinateReferenceSystem) Func {
	return from.To(crs)
}

// ToXYZ converts geographic WGS84 coordinates to geocentric
// coordinates. The incoming parameters are lon, lat, h in units
// of degrees and the resulting parameters are x, y, z in units
// of meters.
func ToXYZ() Func {
	return LonLat().ToXYZ()
}

// ToWebMercator converts geographic WGS84 coordinates to web mercator
// coordinates. The incoming parameters are lon, lat, h in units
// of degrees and the resulting parameters are east, north, h in units
// of meters.
func ToWebMercator() Func {
	return LonLat().ToWebMercator()
}

// ToUTM converts geographic WGS84 coordinates to utm
// coordinates. The incoming parameters are lon, lat, h in units
// of degrees and the resulting parameters are east, north, h in units
// of meters.
func ToUTM(zone float64, northern bool) Func {
	return LonLat().ToUTM(zone, northern)
}

// FromXYZ converts geocentric WGS84 coordinates to geographic
// coordinates. The incoming parameters are x, y, z in units
// of meters and the resulting parameters are lon, lat, h in
// units of degrees..
func FromXYZ() Func {
	return LonLat().FromXYZ()
}

// FromWebMercator converts web mercator WGS84 coordinates to geographic
// coordinates. The incoming parameters are east, north, h in units
// of meters and the resulting parameters are lon, lat, h in
// units of degrees..
func FromWebMercator() Func {
	return LonLat().FromWebMercator()
}

// FromUTM converts utm WGS84 coordinates to geographic
// coordinates. The incoming parameters are east, north, h in units
// of meters and the resulting parameters are lon, lat, h in
// units of degrees..
func FromUTM(zone float64, northern bool) Func {
	return LonLat().FromUTM(zone, northern)
}

// To transforms geographic WGS84 coordinates to another
// CoordinateReferenceSystem. The incoming parameters are
// lon, lat, h in units of degrees.
func To(to CoordinateReferenceSystem) Func {
	return LonLat().To(to)
}

// From transforms coordinates from a CoordinateReferenceSystem
// to geographic WGS84 coordinates. The resulting parameters
// are lon, lat, h in units of degrees.
func From(from CoordinateReferenceSystem) Func {
	return LonLat().From(from)
}
