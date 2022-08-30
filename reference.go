//nolint:varnamelen,nonamedreturns,gomnd
package wgs84

import (
	"errors"
)

// To provides the transformation of WGS84 geographic coordinates to another
// Coordinate Reference System.
func To(to CoordinateReferenceSystem) Func {
	return LonLat().To(to)
}

// From provides the transformation of coordinates from a Coordinate Reference
// System to WGS84 geographic coordinates.
func From(from CoordinateReferenceSystem) Func {
	return LonLat().From(from)
}

// XYZ is a geocentric Coordinate Reference System similar to
// https://epsg.io/4978
func XYZ() GeocentricReferenceSystem {
	return WGS84().XYZ()
}

// LonLat is a geographic Coordinate Reference System similar to
// https://epsg.io/4326
func LonLat() GeographicReferenceSystem {
	return WGS84().LonLat()
}

// WebMercator is a projected Coordinate Reference System similar to
// https://epsg.io/3857
func WebMercator() ProjectedReferenceSystem {
	return WGS84().WebMercator()
}

// UTM represents projected Coordinate Reference System's similar to
// https://epsg.io/32632 or https://epsg.io/32732
func UTM(zone float64, northern bool) ProjectedReferenceSystem {
	northf := 0.0
	if !northern {
		northf = 10000000
	}

	crs := WGS84().TransverseMercator(zone*6-183, 0, 0.9996, 500000, northf)

	crs.Area = AreaFunc(func(lon, lat float64) bool {
		if northern {
			return lon >= zone*6-186 && lon <= zone*6-180 && lat >= 0 && lat <= 84
		}

		return lon >= zone*6-186 && lon <= zone*6-180 && lat <= 0 && lat >= -80
	})

	return crs
}

// ETRS89UTM represents projected Coordinate Reference System's similar to
// https://epsg.io/25832
func ETRS89UTM(zone float64) ProjectedReferenceSystem {
	crs := ETRS89().TransverseMercator(zone*6-183, 0, 0.9996, 500000, 0)
	crs.Area = AreaFunc(func(lon, lat float64) bool {
		return lon >= zone*6-186 && lon <= zone*6-180 && lat >= 0 && lat <= 84
	})

	return crs
}

// ETRS89AustriaLambert represents projected Coordinate Reference System's similar to
// https://epsg.io/3416
func ETRS89AustriaLambert() ProjectedReferenceSystem {
	return ETRS89().LambertConformalConic2SP(13.33333333333333, 47.5, 49, 46, 400000, 400000)
}

func ETRS89LambertAzimuthalEqualArea() ProjectedReferenceSystem {
	return ETRS89().LambertAzimuthalEqualArea(10, 52, 4321000, 3210000)
}

// MGIAustriaLambert represents projected Coordinate Reference System's similar to
// https://epsg.io/31287
func MGIAustriaLambert() ProjectedReferenceSystem {
	return MGI().LambertConformalConic2SP(13.33333333333333, 47.5, 49, 46, 400000, 400000)
}

// MGIAustriaM28 represents projected Coordinate Reference System's similar to
// https://epsg.io/31284
func MGIAustriaM28() ProjectedReferenceSystem {
	return MGI().TransverseMercator(10.33333333333333, 0, 1, 150000, 0)
}

// MGIAustriaM31 represents projected Coordinate Reference System's similar to
// https://epsg.io/31285
func MGIAustriaM31() ProjectedReferenceSystem {
	return MGI().TransverseMercator(13.33333333333333, 0, 1, 450000, 0)
}

// MGIAustriaM34 represents projected Coordinate Reference System's similar to
// https://epsg.io/31286
func MGIAustriaM34() ProjectedReferenceSystem {
	return MGI().TransverseMercator(16.33333333333333, 0, 1, 750000, 0)
}

// MGIAustriaGKM28 represents projected Coordinate Reference System's similar to
// https://epsg.io/31257
func MGIAustriaGKM28() ProjectedReferenceSystem {
	return MGI().TransverseMercator(10.33333333333333, 0, 1, 150000, -5000000)
}

// MGIAustriaGKM31 represents projected Coordinate Reference System's similar to
// https://epsg.io/31258
func MGIAustriaGKM31() ProjectedReferenceSystem {
	return MGI().TransverseMercator(13.33333333333333, 0, 1, 450000, -5000000)
}

// MGIAustriaGKM34 represents projected Coordinate Reference System's similar to
// https://epsg.io/31259
func MGIAustriaGKM34() ProjectedReferenceSystem {
	return MGI().TransverseMercator(16.33333333333333, 0, 1, 750000, -5000000)
}

// OSGB36NationalGrid is a projected Coordinate Reference System similar to
// https://epsg.io/27700
func OSGB36NationalGrid() ProjectedReferenceSystem {
	return OSGB36().TransverseMercator(-2, 49, 0.9996012717, 400000, -100000)
}

// DHDN2001GK represents projected Coordinate Reference System's similar to
// https://epsg.io/31467
func DHDN2001GK(zone float64) ProjectedReferenceSystem {
	crs := DHDN2001().TransverseMercator(zone*3, 0, 1, zone*1000000+500000, 0)
	crs.Area = AreaFunc(func(lon, lat float64) bool {
		return lon >= zone*3-1.5 && lon <= zone*3+1.5 && lat >= 0 && lat <= 84
	})

	return crs
}

// RGF93CC represents projected Coordinate Reference System's similar to
// https://epsg.io/3950
func RGF93CC(lat float64) ProjectedReferenceSystem {
	return RGF93().LambertConformalConic2SP(3, lat, lat-0.75, lat+0.75, 1700000, 2200000+(lat-43)*1000000)
}

// RGF93FranceLambert is a projected Coordinate Reference System similar to
// https://epsg.io/2154
func RGF93FranceLambert() ProjectedReferenceSystem {
	return RGF93().LambertConformalConic2SP(3, 46.5, 49, 44, 700000, 6600000)
}

// NAD83AlabamaEast is a projected Coordinate Reference System similar to
// https://epsg.io/6355
func NAD83AlabamaEast() ProjectedReferenceSystem {
	crs := NAD83().TransverseMercator(-85.83333333333333, 30.5, 0.99996, 200000, 0)
	crs.Area = AreaFunc(func(lon, lat float64) bool {
		return lon >= -86.79 && lon <= -84.89 && lat >= 30.99 && lat <= 35.0
	})

	return crs
}

// NAD83AlabamaWest is a projected Coordinate Reference System similar to
// https://epsg.io/6356
func NAD83AlabamaWest() ProjectedReferenceSystem {
	crs := NAD83().TransverseMercator(-87.5, 30, 0.999933333, 600000, 0)
	crs.Area = AreaFunc(func(lon, lat float64) bool {
		return lon >= -88.48 && lon <= -86.3 && lat >= 30.14 && lat <= 35.02
	})

	return crs
}

// NAD83CaliforniaAlbers is a projected Coordinate Reference System similar to
// https://epsg.io/6414
func NAD83CaliforniaAlbers() ProjectedReferenceSystem {
	crs := NAD83().AlbersEqualAreaConic(34, 40.5, 0, -120, 0, -4000000)
	crs.Area = AreaFunc(func(lon, lat float64) bool {
		return lon >= -124.45 && lon <= -114.12 && lat >= 32.53 && lat <= 42.01
	})

	return crs
}

// GeocentricReferenceSystem represents a geocentric Coordinate Reference System.
type GeocentricReferenceSystem struct {
	Datum Datum
}

// Contains method is the implementation of the Area interface.
func (crs GeocentricReferenceSystem) Contains(lon, lat float64) bool {
	return crs.Datum.Contains(lon, lat)
}

// ToWGS84 method is one method of the CoordinateReferenceSystem interface.
func (crs GeocentricReferenceSystem) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return crs.Datum.Forward(x, y, z)
}

// FromWGS84 method is one method of the CoordinateReferenceSystem interface.
func (crs GeocentricReferenceSystem) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return crs.Datum.Inverse(x0, y0, z0)
}

// To provides the transformation to another CoordinateReferenceSystem.
func (crs GeocentricReferenceSystem) To(to CoordinateReferenceSystem) Func {
	return Transform(crs, to)
}

// SafeTo provides the transformation to another CoordinateReferenceSystem
// with errors.
func (crs GeocentricReferenceSystem) SafeTo(to CoordinateReferenceSystem) SafeFunc {
	return SafeTransform(crs, to)
}

// From provides the transformation from another CoordinateReferenceSystem.
func (crs GeocentricReferenceSystem) From(from CoordinateReferenceSystem) Func {
	return Transform(from, crs)
}

// SafeFrom provides the transformation from another CoordinateReferenceSystem
// with errors.
func (crs GeocentricReferenceSystem) SafeFrom(from CoordinateReferenceSystem) SafeFunc {
	return SafeTransform(from, crs)
}

// GeographicReferenceSystem represents a geographic Coordinate Reference System.
type GeographicReferenceSystem struct {
	Datum Datum
}

// Contains method is the implementation of the Area interface.
func (crs GeographicReferenceSystem) Contains(lon, lat float64) bool {
	return crs.Datum.Contains(lon, lat)
}

// ToWGS84 method is one method of the CoordinateReferenceSystem interface.
func (crs GeographicReferenceSystem) ToWGS84(lon, lat, h float64) (x0, y0, z0 float64) {
	x, y, z := lonLatToXYZ(lon, lat, h, crs.Datum.A(), crs.Datum.Fi())

	return crs.Datum.Forward(x, y, z)
}

// FromWGS84 method is one method of the CoordinateReferenceSystem interface.
func (crs GeographicReferenceSystem) FromWGS84(x0, y0, z0 float64) (lon, lat, h float64) {
	x, y, z := crs.Datum.Inverse(x0, y0, z0)

	return xyzToLonLat(x, y, z, crs.Datum.A(), crs.Datum.Fi())
}

// To provides the transformation to another CoordinateReferenceSystem.
func (crs GeographicReferenceSystem) To(to CoordinateReferenceSystem) Func {
	return Transform(crs, to)
}

// SafeTo provides the transformation to another CoordinateReferenceSystem
// with errors.
func (crs GeographicReferenceSystem) SafeTo(to CoordinateReferenceSystem) SafeFunc {
	return SafeTransform(crs, to)
}

// From provides the transformation from another CoordinateReferenceSystem.
func (crs GeographicReferenceSystem) From(from CoordinateReferenceSystem) Func {
	return Transform(from, crs)
}

// SafeFrom provides the transformation from another CoordinateReferenceSystem
// with errors.
func (crs GeographicReferenceSystem) SafeFrom(from CoordinateReferenceSystem) SafeFunc {
	return SafeTransform(from, crs)
}

// ProjectedReferenceSystem represents a projected Coordinate Reference System.
type ProjectedReferenceSystem struct {
	Datum      Datum
	Projection Projection
	Area       Area
}

// Contains method is the implementation of the Area interface.
func (crs ProjectedReferenceSystem) Contains(lon, lat float64) bool {
	return crs.Datum.Contains(lon, lat) && (crs.Area == nil || crs.Area.Contains(lon, lat))
}

// ToWGS84 method is one method of the CoordinateReferenceSystem interface.
func (crs ProjectedReferenceSystem) ToWGS84(east, north, h float64) (x0, y0, z0 float64) {
	if crs.Projection == nil {
		return crs.Datum.WebMercator().ToWGS84(east, north, h)
	}

	lon, lat := crs.Projection.ToLonLat(east, north, crs.Datum)
	x, y, z := lonLatToXYZ(lon, lat, h, crs.Datum.A(), crs.Datum.Fi())

	return crs.Datum.Forward(x, y, z)
}

// FromWGS84 method is one method of the CoordinateReferenceSystem interface.
func (crs ProjectedReferenceSystem) FromWGS84(x0, y0, z0 float64) (east, north, h float64) {
	if crs.Projection == nil {
		return crs.Datum.WebMercator().FromWGS84(x0, y0, z0)
	}

	x, y, z := crs.Datum.Inverse(x0, y0, z0)
	lon, lat, h := xyzToLonLat(x, y, z, crs.Datum.A(), crs.Datum.Fi())
	east, north = crs.Projection.FromLonLat(lon, lat, crs.Datum)

	return east, north, h
}

// To provides the transformation to another CoordinateReferenceSystem.
func (crs ProjectedReferenceSystem) To(to CoordinateReferenceSystem) Func {
	return Transform(crs, to)
}

// SafeTo provides the transformation to another CoordinateReferenceSystem
// with errors.
func (crs ProjectedReferenceSystem) SafeTo(to CoordinateReferenceSystem) SafeFunc {
	return SafeTransform(crs, to)
}

// From provides the transformation from another CoordinateReferenceSystem.
func (crs ProjectedReferenceSystem) From(from CoordinateReferenceSystem) Func {
	return Transform(from, crs)
}

// SafeFrom provides the transformation from another CoordinateReferenceSystem
// with errors.
func (crs ProjectedReferenceSystem) SafeFrom(from CoordinateReferenceSystem) SafeFunc {
	return SafeTransform(from, crs)
}

// Transform provides a transformation between CoordinateReferenceSystems.
func Transform(from, to CoordinateReferenceSystem) Func {
	return func(a, b, c float64) (a2, b2, c2 float64) {
		if from != nil {
			a, b, c = from.ToWGS84(a, b, c)
		}

		if to != nil {
			a, b, c = to.FromWGS84(a, b, c)
		}

		return a, b, c
	}
}

var (
	// ErrNoCoordinateReferenceSystem is a nil CoordinateReferenceSystem warning.
	ErrNoCoordinateReferenceSystem = errors.New("crs not specified")
	// ErrOutOfBounds is a transformation out of the Area interface boundings.
	ErrOutOfBounds = errors.New("coordinate is out of bounds")
)

// SafeTransform provides a transformation between CoordinateReferenceSystems
// with errors.
func SafeTransform(from, to CoordinateReferenceSystem) SafeFunc {
	return func(a, b, c float64) (float64, float64, float64, error) {
		if from == nil || to == nil {
			return 0, 0, 0, ErrNoCoordinateReferenceSystem
		}

		a, b, c = from.ToWGS84(a, b, c)

		lon, lat, _ := xyzToLonLat(a, b, c, A, Fi)
		if !from.Contains(lon, lat) || !to.Contains(lon, lat) {
			return 0, 0, 0, ErrOutOfBounds
		}

		a, b, c = to.FromWGS84(a, b, c)

		return a, b, c, nil
	}
}
