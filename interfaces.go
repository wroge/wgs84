package wgs84

// GeodeticSpheroid interface is used by Spheroid, Datum, ETRS89, LonLat, ...
type GeodeticSpheroid interface {
	MajorAxis() float64
	InverseFlattening() float64
}

// Transformation interface is used by Helmert, Datum, ETRS89, LonLat, ...
type Transformation interface {
	ToWGS84(x, y, z float64) (x0, y0, z0 float64)
	FromWGS84(x0, y0, z0 float64) (x, y, z float64)
}

// GeodeticDatum composes the GeodeticSpheroid, Transformation and Area
// interface.
type GeodeticDatum interface {
	GeodeticSpheroid
	Transformation
	Area
}

// CoordinateSystem is used by XYZ, LonLat, WebMercator, ...
type CoordinateSystem interface {
	ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64)
	FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64)
}

// CoordinateProjection interface allows easy extension of the package.
type CoordinateProjection interface {
	ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64)
	FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64)
}

// CoordinateReferenceSystem composes CoordinateSystem and GeodeticDatum
// interface.
type CoordinateReferenceSystem interface {
	CoordinateSystem
	GeodeticDatum
}

// Area interface used by Repository to detect available EPSG-Codes at a
// specific location.
type Area interface {
	Contains(lon, lat float64) bool
}
