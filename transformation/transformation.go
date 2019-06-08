package transformation

import "math"

func WGS84() Transformation {
	return Transformation{
		toWGS84: func(x, y, z float64) (x0, y0, z0 float64) {
			return x, y, z
		},
		fromWGS84: func(x0, y0, z0 float64) (x, y, z float64) {
			return x0, y0, z0
		},
	}
}

type Transformation struct {
	toWGS84   func(x, y, z float64) (x0, y0, z0 float64)
	fromWGS84 func(x0, y0, z0 float64) (x, y, z float64)
}

func (tra Transformation) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	if tra.toWGS84 == nil {
		return x, y, z
	}
	return tra.toWGS84(x, y, z)
}

func (tra Transformation) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	if tra.fromWGS84 == nil {
		return x0, y0, z0
	}
	return tra.fromWGS84(x0, y0, z0)
}

func (t Transformation) GeocentricTranslation(tx, ty, tz float64) Transformation {
	return Transformation{
		toWGS84: func(x, y, z float64) (x0, y0, z0 float64) {
			x, y, z = t.ToWGS84(x, y, z)
			return x + tx, y + ty, z + tz
		},
		fromWGS84: func(x0, y0, z0 float64) (x, y, z float64) {
			x0, y0, z0 = x0-tx, y0-ty, z0-tz
			return t.FromWGS84(x0, y0, z0)
		},
	}
}

func (t Transformation) Helmert(tx, ty, tz, rx, ry, rz, ds float64) Transformation {
	return Transformation{
		toWGS84: func(x, y, z float64) (x0, y0, z0 float64) {
			x, y, z = t.ToWGS84(x, y, z)
			return helmert(x, y, z, tx, ty, tz, rx, ry, rz, ds)
		},
		fromWGS84: func(x0, y0, z0 float64) (x, y, z float64) {
			x0, y0, z0 = helmert(x0, y0, z0, -tx, -ty, -tz, -rx, -ry, -rz, -ds)
			return t.FromWGS84(x0, y0, z0)
		},
	}
}

func helmert(x, y, z, tx, ty, tz, rx, ry, rz, ds float64) (x0, y0, z0 float64) {
	asec := math.Pi / 648000
	ppm := 0.000001
	x0 = (1+ds*ppm)*(x+z*ry*asec-y*rz*asec) + tx
	y0 = (1+ds*ppm)*(y+x*rz*asec-z*rx*asec) + ty
	z0 = (1+ds*ppm)*(z+y*rx*asec-x*ry*asec) + tz
	return
}
