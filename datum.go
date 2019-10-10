package wgs84

import "math"

// Datum provides the GeodeticDatum interface with customizable parameters.
type Datum struct {
	GeodeticSpheroid GeodeticSpheroid
	Transformation   Transformation
	Area             Area
}

// Contains is the implementation of the Area interface.
//
// By default Contains is true between -180 to 180 degrees longitude and
// -90 and 90 degrees latitude. The Area kann be limited by Datum.Area.
func (d Datum) Contains(lon, lat float64) bool {
	if d.Area != nil && !d.Area.Contains(lon, lat) {
		return false
	}
	if math.Abs(lat) > 180 || math.Abs(lat) > 90 {
		return false
	}
	return true
}

// MajorAxis is a method for the implementation of the GeodeticSpheroid
// interface. By default the MajorAxis of WGS84 is returned.
func (d Datum) MajorAxis() float64 {
	return spheroid(d.GeodeticSpheroid).MajorAxis()
}

// InverseFlattening is a method for the implementation of the
// GeodeticSpheroid interface. By default the InverseFlattening of WGS84 is
// returned.
func (d Datum) InverseFlattening() float64 {
	return spheroid(d.GeodeticSpheroid).InverseFlattening()
}

// ToWGS84 is a method for the implementation of the Transformation
// interface. By default x, y and z are returned.
func (d Datum) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(d.Transformation, x, y, z)
}

// FromWGS84 is a method for the implementation of the Transformation
// interface. By default x0, y0 and z0 are returned.
func (d Datum) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(d.Transformation, x0, y0, z0)
}

// ETRS89 provides the GeodeticDatum interface for the European Terrestrial
// Reference System 1989.
type ETRS89 struct {
	Area Area
}

// Contains is the implementation of the Area interface.
//
// By default Contains is true between -180 to 180 degrees longitude and
// -90 and 90 degrees latitude. The Area kann be limited by ETRS89.Area.
func (d ETRS89) Contains(lon, lat float64) bool {
	if d.Area != nil && !d.Area.Contains(lon, lat) {
		return false
	}
	if lon < -16.1 || lon > 40.18 || lat < 32.88 || lat > 84.17 {
		return false
	}
	return true
}

// MajorAxis is a method for the implementation of the GeodeticSpheroid
// interface. By default the MajorAxis of WGS84 is returned.
func (d ETRS89) MajorAxis() float64 {
	return GRS80().MajorAxis()
}

// InverseFlattening is a method for the implementation of the
// GeodeticSpheroid interface. By default the InverseFlattening of WGS84 is
// returned.
func (d ETRS89) InverseFlattening() float64 {
	return GRS80().InverseFlattening()
}

// ToWGS84 is a method for the implementation of the Transformation
// interface. By default x, y and z are returned.
func (d ETRS89) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return x, y, z
}

// FromWGS84 is a method for the implementation of the Transformation
// interface. By default x0, y0 and z0 are returned.
func (d ETRS89) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return x0, y0, z0
}

// LonLat implements a CoordinateReferenceSystem similar to the EPSG-Code
// 4258.
func (d ETRS89) LonLat() LonLat {
	return LonLat{
		GeodeticDatum: d,
	}
}

// UTM implements CoordinateReferenceSystems similar to the EPSG-Codes
// 25828 - 25838.
func (d ETRS89) UTM(zone float64) TransverseMercator {
	return TransverseMercator{
		GeodeticDatum: d,
		Lonf:          zone*6 - 183,
		Latf:          0,
		Scale:         0.9996,
		Eastf:         500000,
		Northf:        0,
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < zone*6-186 || lon > zone*6-180 {
				return false
			}
			return true
		}),
	}
}

// OSGB36 provides the GeodeticDatum interface for the Ordnance Survey Great
// Britain 1936 Datum.
type OSGB36 struct {
	Area Area
}

// Contains is the implementation of the Area interface.
//
// By default Contains is true between -180 to 180 degrees longitude and
// -90 and 90 degrees latitude. The Area kann be limited by OSGB36.Area.
func (d OSGB36) Contains(lon, lat float64) bool {
	if d.Area != nil && !d.Area.Contains(lon, lat) {
		return false
	}
	if lon < -8.82 || lon > 1.92 || lat < 49.79 || lat > 60.94 {
		return false
	}
	return true
}

func (d OSGB36) transformation() Helmert {
	return Helmert{
		Tx: 446.448,
		Ty: -125.157,
		Tz: 542.06,
		Rx: 0.15,
		Ry: 0.247,
		Rz: 0.842,
		Ds: -20.489,
	}
}

// MajorAxis is a method for the implementation of the GeodeticSpheroid
// interface. By default the MajorAxis of WGS84 is returned.
func (d OSGB36) MajorAxis() float64 {
	return Airy().MajorAxis()
}

// InverseFlattening is a method for the implementation of the
// GeodeticSpheroid interface. By default the InverseFlattening of WGS84 is
// returned.
func (d OSGB36) InverseFlattening() float64 {
	return Airy().InverseFlattening()
}

// ToWGS84 is a method for the implementation of the Transformation
// interface. By default x, y and z are returned.
func (d OSGB36) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return d.transformation().ToWGS84(x, y, z)
}

// FromWGS84 is a method for the implementation of the Transformation
// interface. By default x0, y0 and z0 are returned.
func (d OSGB36) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return d.transformation().FromWGS84(x0, y0, z0)
}

// LonLat implements a CoordinateReferenceSystem similar to the EPSG-Code
// 4277.
func (d OSGB36) LonLat() LonLat {
	return LonLat{
		GeodeticDatum: d,
	}
}

// NationalGrid implements a CoordinateReferenceSystem similar to the
// EPSG-Code 27700.
func (d OSGB36) NationalGrid() TransverseMercator {
	return TransverseMercator{
		GeodeticDatum: d,
		Lonf:          -2,
		Latf:          49,
		Scale:         0.9996012717,
		Eastf:         400000,
		Northf:        -100000,
	}
}

// DHDN2001 provides the GeodeticDatum interface for the "Deutsche
// Hauptdreiecksnetz" 2001 Datum in Germany.
type DHDN2001 struct {
	Area Area
}

// Contains is the implementation of the Area interface.
//
// By default Contains is true between -180 to 180 degrees longitude and
// -90 and 90 degrees latitude. The Area kann be limited by DHDN2001.Area.
func (d DHDN2001) Contains(lon, lat float64) bool {
	if d.Area != nil && !d.Area.Contains(lon, lat) {
		return false
	}
	if lon < 5.87 || lon > 13.84 || lat < 47.27 || lat > 55.09 {
		return false
	}
	return true
}

func (d DHDN2001) transformation() Helmert {
	return Helmert{
		Tx: 598.1,
		Ty: 73.7,
		Tz: 418.2,
		Rx: 0.202,
		Ry: 0.045,
		Rz: -2.455,
		Ds: 6.7,
	}
}

// MajorAxis is a method for the implementation of the GeodeticSpheroid
// interface. By default the MajorAxis of WGS84 is returned.
func (d DHDN2001) MajorAxis() float64 {
	return Bessel().MajorAxis()
}

// InverseFlattening is a method for the implementation of the
// GeodeticSpheroid interface. By default the InverseFlattening of WGS84 is
// returned.
func (d DHDN2001) InverseFlattening() float64 {
	return Bessel().InverseFlattening()
}

// ToWGS84 is a method for the implementation of the Transformation
// interface. By default x, y and z are returned.
func (d DHDN2001) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return d.transformation().ToWGS84(x, y, z)
}

// FromWGS84 is a method for the implementation of the Transformation
// interface. By default x0, y0 and z0 are returned.
func (d DHDN2001) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return d.transformation().FromWGS84(x0, y0, z0)
}

// LonLat implements a CoordinateReferenceSystem similar to the
// EPSG-Code 4314.
func (d DHDN2001) LonLat() LonLat {
	return LonLat{
		GeodeticDatum: d,
	}
}

// GK implements CoordinateReferenceSystems similar to the EPSG-Codes
// 31466- 31469.
func (d DHDN2001) GK(zone float64) TransverseMercator {
	return TransverseMercator{
		GeodeticDatum: d,
		Lonf:          zone * 3,
		Latf:          0,
		Scale:         1,
		Eastf:         zone*1000000 + 500000,
		Northf:        0,
	}
}

// RGF93 provides the GeodeticDatum interface for the "Réseau géodésique
// français" 1993 Datum in France.
type RGF93 struct {
	Area Area
}

// Contains is the implementation of the Area interface.
//
// By default Contains is true between -180 to 180 degrees longitude and
// -90 and 90 degrees latitude. The Area kann be limited by RGF93.Area.
func (d RGF93) Contains(lon, lat float64) bool {
	if d.Area != nil && !d.Area.Contains(lon, lat) {
		return false
	}
	if lon < -9.86 || lon > 10.38 || lat < 41.15 || lat > 51.56 {
		return false
	}
	return true
}

// MajorAxis is a method for the implementation of the GeodeticSpheroid
// interface. By default the MajorAxis of WGS84 is returned.
func (d RGF93) MajorAxis() float64 {
	return GRS80().MajorAxis()
}

// InverseFlattening is a method for the implementation of the
// GeodeticSpheroid interface. By default the InverseFlattening of WGS84 is
// returned.
func (d RGF93) InverseFlattening() float64 {
	return GRS80().InverseFlattening()
}

// ToWGS84 is a method for the implementation of the Transformation
// interface. By default x, y and z are returned.
func (d RGF93) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return x, y, z
}

// FromWGS84 is a method for the implementation of the Transformation
// interface. By default x0, y0 and z0 are returned.
func (d RGF93) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return x0, y0, z0
}

// LonLat implements a CoordinateReferenceSystem similar to the
// EPSG-Code 4314.
func (d RGF93) LonLat() LonLat {
	return LonLat{
		GeodeticDatum: d,
	}
}

// CC implements CoordinateReferenceSystems similar to the EPSG-Codes
// 3942- 3950.
func (d RGF93) CC(lat float64) LambertConformalConic2SP {
	return LambertConformalConic2SP{
		GeodeticDatum: d,
		Lonf:          3,
		Latf:          lat,
		Lat1:          lat - 0.75,
		Lat2:          lat + 0.75,
		Eastf:         1700000,
		Northf:        2200000 + (lat-43)*1000000,
	}
}

// FranceLambert implements a CoordinateReferenceSystem similar to the
// EPSG-Code 2154.
func (d RGF93) FranceLambert() LambertConformalConic2SP {
	return LambertConformalConic2SP{
		GeodeticDatum: d,
		Lonf:          3,
		Latf:          46.5,
		Lat1:          49,
		Lat2:          44,
		Eastf:         700000,
		Northf:        6600000,
	}
}
