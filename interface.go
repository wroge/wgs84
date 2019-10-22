package wgs84

// Spheroid interface represents an ellipsoid of revolution used by
// geodetic datums.
//
// It is specified by the major axis and the inverse flattening.
type Spheroid interface {
	A() float64
	Fi() float64
}

// Transformation interface represents the transformation of geocentric
// coordinates To and From WGS84.
//
// The Forward method is used to transform coordinates to the WGS84 System..
type Transformation interface {
	Forward(x, y, z float64) (x0, y0, z0 float64)
	Inverse(x0, y0, z0 float64) (x, y, z float64)
}

// Projection interface is used by the several Projected Coordinate
// Reference System's in this package.
//
// It is easy to expand the package with additional Projections
// through this interface.
type Projection interface {
	ToLonLat(east, north float64, s Spheroid) (lon, lat float64)
	FromLonLat(lon, lat float64, s Spheroid) (east, north float64)
}

// CoordinateReferenceSystem is the core interface of this package.
//
// It is a coordinate reference system that precisely locates a point
// on the earth's surface.
type CoordinateReferenceSystem interface {
	ToWGS84(a, b, c float64) (x0, y0, z0 float64)
	FromWGS84(x0, y0, z0 float64) (a, b, c float64)
	Area
}

// Area interface is used to describe the Bounding Box of a Coordinate
// Reference System.
//
// It is implemented by the AreaFunc.
type Area interface {
	Contains(lon, lat float64) bool
}
