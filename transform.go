package wgs84

import (
	"math"
)

// Transform is the core utility of this package. It is used by the To
// methods of each CoordinateReferenceSystem in this package to transform
// coordinates from one CoordinateReferenceSystem to another.
func Transform(from, to CoordinateReferenceSystem) Func {
	return func(a, b, c float64) (a2, b2, c2 float64) {
		x, y, z := from.ToXYZ(a, b, c, from)
		x0, y0, z0 := from.ToWGS84(x, y, z)
		x2, y2, z2 := to.FromWGS84(x0, y0, z0)
		return to.FromXYZ(x2, y2, z2, to)
	}
}

type Func func(a, b, c float64) (a2, b2, c2 float64)

func (f Func) Round(precision float64) Func {
	round := func(v float64) float64 {
		r := math.Round(v*math.Pow(10, precision)) / math.Pow(10, precision)
		if r == -0 {
			return 0
		}
		return r
	}
	return func(a, b, c float64) (a2, b2, c2 float64) {
		a, b, c = f(a, b, c)
		return round(a), round(b), round(c)
	}
}
