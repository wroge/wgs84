package wgs84

type Spheroid interface {
	A() float64
	Fi() float64
}

type Transformation interface {
	Forward(x, y, z float64) (x0, y0, z0 float64)
	Inverse(x0, y0, z0 float64) (x, y, z float64)
}

type Projection interface {
	ToLonLat(east, north float64, s Spheroid) (lon, lat float64)
	FromLonLat(lon, lat float64, s Spheroid) (east, north float64)
}

type CoordinateReferenceSystem interface {
	ToWGS84(a, b, c float64) (x0, y0, z0 float64)
	FromWGS84(x0, y0, z0 float64) (a, b, c float64)
	Area
}

type Area interface {
	Contains(lon, lat float64) bool
}
