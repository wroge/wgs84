package wgs84

type GeodeticSpheroid interface {
	MajorAxis() float64
	InverseFlattening() float64
}

type Transformation interface {
	ToWGS84(x, y, z float64) (x0, y0, z0 float64)
	FromWGS84(x0, y0, z0 float64) (x, y, z float64)
}

type GeodeticDatum interface {
	GeodeticSpheroid
	Transformation
}

type CoordinateSystem interface {
	ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64)
	FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64)
}

type CoordinateProjection interface {
	ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64)
	FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64)
}

type CoordinateReferenceSystem interface {
	GeodeticDatum
	CoordinateSystem
}
