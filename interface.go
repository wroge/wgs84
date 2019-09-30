package wgs84

import "github.com/wroge/wgs84/spheroid"

type Spheroid = interface {
	A() float64
	Fi() float64
}

type System interface {
	ToXYZ(a, b, c float64, s Spheroid) (x, y, z float64)
	FromXYZ(x, y, z float64, s Spheroid) (a, b, c float64)
}

type Transformation interface {
	ToWGS84(x, y, z float64) (x0, y0, z0 float64)
	FromWGS84(x0, y0, z0 float64) (x, y, z float64)
}

type CoordinateReferenceSystem struct {
	Spheroid       Spheroid
	Transformation Transformation
	System         System
}

// Empty CoordinateReferenceSystems are handled as WGS84 XYZ (EPSG-Code 4978)
func (crs CoordinateReferenceSystem) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	if crs.Spheroid == nil {
		crs.Spheroid = spheroid.WGS84()
	}
	if to.Spheroid == nil {
		to.Spheroid = spheroid.WGS84()
	}
	return func(a, b, c float64) (a2, b2, c2 float64) {
		if crs.System != nil {
			a, b, c = crs.System.ToXYZ(a, b, c, crs.Spheroid)
		}
		if crs.Transformation != nil {
			a, b, c = crs.Transformation.ToWGS84(a, b, c)
		}
		if to.Transformation != nil {
			a, b, c = to.Transformation.FromWGS84(a, b, c)
		}
		if to.System != nil {
			a, b, c = to.System.FromXYZ(a, b, c, to.Spheroid)
		}
		return a, b, c
	}
}
