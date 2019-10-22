package wgs84

import (
	"math"
)

// Helmert provides a Datum specified through the major axis and the
// inverse flattening of a spheroid and the 7 parameters of a Helmert-
// Transformation.
func Helmert(a, fi, tx, ty, tz, rx, ry, rz, ds float64) Datum {
	return Datum{
		Spheroid: spheroid{a: a, fi: fi},
		Transformation: helmert{
			tx: tx,
			ty: ty,
			tz: tz,
			rx: rx,
			ry: ry,
			rz: rz,
			ds: ds,
		},
	}
}

// WGS84 provides a Datum similar to the World Geodetic System 1984.
//
// It's based on the WGS84 Spheroid.
//
// It is used worldwide.
func WGS84() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6378137,
			fi: 298.257223563,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if math.Abs(lon) > 180 || math.Abs(lat) > 90 {
				return false
			}
			return true
		}),
	}
}

// ETRS89 provides a Datum similar to the European Terrestrial Reference
// System 1989.
//
// It's based on the GRS80 Spheroid.
//
// It is used in Europe.
func ETRS89() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6378137,
			fi: 298.257222101,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < -16.1 || lon > 40.18 || lat < 32.88 || lat > 84.17 {
				return false
			}
			return true
		}),
	}
}

// OSGB36 provides a Datum similar to the Ordnance Survey Great Britain 1936.
//
// It's based on the Airy Spheroid and a 7-parameter-Helmert-Transformation
// with the parameters: 446.448,-125.157,542.06,0.15,0.247,0.842,-20.489.
//
// https://epsg.io/1314
//
// It is used in Great Britain.
func OSGB36() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6377563.396,
			fi: 299.3249646,
		},
		Transformation: helmert{
			tx: 446.448,
			ty: -125.157,
			tz: 542.06,
			rx: 0.15,
			ry: 0.247,
			rz: 0.842,
			ds: -20.489,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < -8.82 || lon > 1.92 || lat < 49.79 || lat > 60.94 {
				return false
			}
			return true
		}),
	}
}

// DHDN2001 provides a Datum similar to the Deutsches Hauptdreiecksnetz 2001.
//
// It's based on the Bessel Spheroid and a 7-parameter-Helmert-Transformation
// with the parameters: 598.1, 73.7, 418.2, 0.202, 0.045, -2.455, 6.7.
//
// https://epsg.io/1776
//
// It is used in Germay.
func DHDN2001() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6377397.155,
			fi: 299.1528128,
		},
		Transformation: helmert{
			tx: 598.1,
			ty: 73.7,
			tz: 418.2,
			rx: 0.202,
			ry: 0.045,
			rz: -2.455,
			ds: 6.7,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < 5.87 || lon > 13.84 || lat < 47.27 || lat > 55.09 {
				return false
			}
			return true
		}),
	}
}

// RGF93 provides a Datum similar to the Réseau géodésique français 1993.
//
// It's based on the GRS80 Spheroid.
//
// It is used in France.
func RGF93() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6378137,
			fi: 298.257222101,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < -9.86 || lon > 10.38 || lat < 41.15 || lat > 51.56 {
				return false
			}
			return true
		}),
	}
}

// NAD83 provides a Datum similar to the North American Datum 1983.
//
// It's based on the GRS80 Spheroid.
//
// It is used in North-America.
func NAD83() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6378137,
			fi: 298.257222101,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < -172.54 || lon > -47.74 || lat < 23.81 || lat > 86.46 {
				return false
			}
			return true
		}),
	}
}

// Datum represents a Geodetic Datum like WGS84, ETRS89 or NAD83.
//
// It implements the Spheroid, Transformation and Area interface.
//
// By default it behaves like a WGS84 Datum.
type Datum struct {
	Spheroid       Spheroid
	Transformation Transformation
	Area           Area
}

// Contains method is the implementation of the Area interface.
//
// Returns false for latitudes with an absolute over 90 and longitudes over 180.
//
// Returns true if nil.
func (d Datum) Contains(lon, lat float64) bool {
	if math.Abs(lon) > 180 || math.Abs(lat) > 90 {
		return false
	}
	if d.Area != nil {
		return d.Area.Contains(lon, lat)
	}
	return true
}

// A is one implementation of the Spheroid interface.
//
// Returns the major axis of the implemented Spheroid.
//
// If nil it returns the major axis of the WGS84 Spheroid.
func (d Datum) A() float64 {
	if d.Spheroid == nil {
		return 6378137
	}
	return d.Spheroid.A()
}

// Fi is one implementation of the Spheroid interface.
//
// Returns the inverse flattening of the implemented Spheroid.
//
// If nil it returns the inverse flattening of the WGS84 Spheroid.
func (d Datum) Fi() float64 {
	if d.Spheroid == nil {
		return 298.257223563
	}
	return d.Spheroid.Fi()
}

// Forward is one implementation of the Tranformation interface.
//
// Returns the Forward transformation of the implemented Transformation.
//
// Returns x, y, z if nil.
func (d Datum) Forward(x, y, z float64) (x0, y0, z0 float64) {
	if d.Transformation == nil {
		return x, y, z
	}
	return d.Transformation.Forward(x, y, z)
}

// Inverse is one implementation of the Tranformation interface.
//
// Returns the Inverse transformation of the implemented Transformation.
//
// Returns x0, y0, z0 if nil.
func (d Datum) Inverse(x0, y0, z0 float64) (x, y, z float64) {
	if d.Transformation == nil {
		return x0, y0, z0
	}
	return d.Transformation.Inverse(x0, y0, z0)
}

// XYZ is a geocentric Coordinate Reference System.
func (d Datum) XYZ() GeocentricReferenceSystem {
	return GeocentricReferenceSystem{
		Datum: d,
	}
}

// LonLat is a geographic Coordinate Reference System.
func (d Datum) LonLat() GeographicReferenceSystem {
	return GeographicReferenceSystem{
		Datum: d,
	}
}

// WebMercator is a projected Coordinate Reference System.
func (d Datum) WebMercator() ProjectedReferenceSystem {
	return ProjectedReferenceSystem{
		Datum:      d,
		Projection: webMercator{},
	}
}

// TransverseMercator is a projected Coordinate Reference System.
func (d Datum) TransverseMercator(lonf, latf, scale, eastf, northf float64) ProjectedReferenceSystem {
	return ProjectedReferenceSystem{
		Datum: d,
		Projection: transverseMercator{
			lonf:   lonf,
			latf:   latf,
			scale:  scale,
			eastf:  eastf,
			northf: northf,
		},
	}
}

// LambertConformalConic2SP is a projected Coordinate Reference System.
func (d Datum) LambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf float64) ProjectedReferenceSystem {
	return ProjectedReferenceSystem{
		Datum: d,
		Projection: lambertConformalConic2SP{
			lonf:   lonf,
			latf:   latf,
			lat1:   lat1,
			lat2:   lat2,
			eastf:  eastf,
			northf: northf,
		},
	}
}

// AlbersEqualAreaConic is a projected Coordinate Reference System.
func (d Datum) AlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf float64) ProjectedReferenceSystem {
	return ProjectedReferenceSystem{
		Datum: d,
		Projection: albersEqualAreaConic{
			lonf:   lonf,
			latf:   latf,
			lat1:   lat1,
			lat2:   lat2,
			eastf:  eastf,
			northf: northf,
		},
	}
}
