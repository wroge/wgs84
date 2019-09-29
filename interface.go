package wgs84

type Spheroid interface {
	A() float64
	Fi() float64
}

type CoordinateSystem interface {
	ToXYZ(a, b, c float64, s Spheroid) (x, y, z float64)
	FromXYZ(x, y, z float64, s Spheroid) (a, b, c float64)
}

type Transformation interface {
	ToWGS84(x, y, z float64) (x0, y0, z0 float64)
	FromWGS84(x0, y0, z0 float64) (x, y, z float64)
}
