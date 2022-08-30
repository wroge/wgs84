//nolint:varnamelen,nonamedreturns,gomnd,exhaustivestruct,exhaustruct
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
			a:  A,
			fi: Fi,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			return math.Abs(lon) <= 180 && math.Abs(lat) <= 90
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
		Spheroid: GRS80{},
		Area: AreaFunc(func(lon, lat float64) bool {
			return lon >= -16.1 && lon <= 40.18 && lat >= 32.88 && lat <= 84.17
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
		Spheroid: Airy{},
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
			return lon >= -8.82 && lon <= 1.92 && lat >= 49.79 && lat <= 60.94
		}),
	}
}

// MGI provides a Datum similar to the Militar-Geographische Institut.
//
// It's based on the Bessel Spheroid and a 7-parameter-Helmert-Transformation
// with the parameters: 577.326,90.129,463.919,5.137,1.474,5.297,2.4232.
//
// https://epsg.io/1618
//
// It is used in Austria.
func MGI() Datum {
	return Datum{
		Spheroid: Bessel{},
		Transformation: helmert{
			tx: 577.326,
			ty: 90.129,
			tz: 463.919,
			rx: 5.137,
			ry: 1.474,
			rz: 5.297,
			ds: 2.4232,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			return lon >= 9.53 && lon <= 17.17 && lat >= 46.4 && lat <= 49.02
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
		Spheroid: Bessel{},
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
			return lon >= 5.87 && lon <= 13.84 && lat >= 47.27 && lat <= 55.09
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
		Spheroid: GRS80{},
		Area: AreaFunc(func(lon, lat float64) bool {
			return lon >= -9.86 && lon <= 10.38 && lat >= 41.15 && lat <= 51.56
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
		Spheroid: GRS80{},
		Area: AreaFunc(func(lon, lat float64) bool {
			return lon >= -172.54 && lon <= -47.74 && lat >= 23.81 && lat <= 86.46
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
	return math.Abs(lon) <= 180 && math.Abs(lat) <= 90 && d.Area != nil && d.Area.Contains(lon, lat)
}

// A returns the major axis of the implemented Spheroid.
//
// If nil it returns the major axis of the WGS84 Spheroid.
func (d Datum) A() float64 {
	if d.Spheroid == nil {
		return A
	}

	return d.Spheroid.A()
}

// Fi returns the inverse flattening of the implemented Spheroid.
//
// If nil it returns the inverse flattening of the WGS84 Spheroid.
func (d Datum) Fi() float64 {
	if d.Spheroid == nil {
		return Fi
	}

	return d.Spheroid.Fi()
}

// Forward transforms geocentric coordinates to WGS84.
//
// Returns x, y, z if nil.
func (d Datum) Forward(x, y, z float64) (x0, y0, z0 float64) {
	if d.Transformation == nil {
		return x, y, z
	}

	return d.Transformation.Forward(x, y, z)
}

// Inverse transforms geocentric coordinates from WGS84.
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

func (d Datum) LambertAzimuthalEqualArea(lonf, latf, eastf, northf float64) ProjectedReferenceSystem {
	return ProjectedReferenceSystem{
		Datum: d,
		Projection: lambertAzimuthalEqualArea{
			latf:   latf,
			lonf:   lonf,
			eastf:  eastf,
			northf: northf,
		},
	}
}
