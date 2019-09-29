package wgs84

import "math"

type transformation struct {
	toWGS84   func(x, y, z float64) (x0, y0, z0 float64)
	fromWGS84 func(x0, y0, z0 float64) (x, y, z float64)
}

func (t transformation) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	if t.toWGS84 == nil {
		return x, y, z
	}
	return t.toWGS84(x, y, z)
}

func (t transformation) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	if t.fromWGS84 == nil {
		return x0, y0, z0
	}
	return t.fromWGS84(x0, y0, z0)
}

func geocentricTranslation(t Transformation, tx, ty, tz float64) transformation {
	return transformation{
		toWGS84: func(x, y, z float64) (x0, y0, z0 float64) {
			if t != nil {
				x, y, z = t.ToWGS84(x, y, z)
			}
			return x + tx, y + ty, z + tz
		},
		fromWGS84: func(x0, y0, z0 float64) (x, y, z float64) {
			x0, y0, z0 = x0-tx, y0-ty, z0-tz
			if t == nil {
				return x0, y0, z0
			}
			return t.FromWGS84(x0, y0, z0)
		},
	}
}

func helmert(t Transformation, tx, ty, tz, rx, ry, rz, ds float64) transformation {
	h := func(x, y, z, tx, ty, tz, rx, ry, rz, ds float64) (x0, y0, z0 float64) {
		asec := math.Pi / 648000
		ppm := 0.000001
		x0 = (1+ds*ppm)*(x+z*ry*asec-y*rz*asec) + tx
		y0 = (1+ds*ppm)*(y+x*rz*asec-z*rx*asec) + ty
		z0 = (1+ds*ppm)*(z+y*rx*asec-x*ry*asec) + tz
		return
	}
	return transformation{
		toWGS84: func(x, y, z float64) (x0, y0, z0 float64) {
			if t != nil {
				x, y, z = t.ToWGS84(x, y, z)
			}
			return h(x, y, z, tx, ty, tz, rx, ry, rz, ds)
		},
		fromWGS84: func(x0, y0, z0 float64) (x, y, z float64) {
			x0, y0, z0 = h(x0, y0, z0, -tx, -ty, -tz, -rx, -ry, -rz, -ds)
			if t == nil {
				return x0, y0, z0
			}
			return t.FromWGS84(x0, y0, z0)
		},
	}
}
