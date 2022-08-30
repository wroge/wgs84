//nolint:nonamedreturns,gomnd
package wgs84

import "math"

// Func is returned by the To methods of the CoordinateReferenceSystem's
// in this package and the Transform function.
type Func func(a, b, c float64) (a2, b2, c2 float64)

func round(v, precision float64) float64 {
	r := math.Round(v*math.Pow(10, precision)) / math.Pow(10, precision)
	if r == -0 {
		return 0
	}

	return r
}

// Round can round the resulting values to a specific precision.
func (f Func) Round(precision float64) Func {
	return func(a, b, c float64) (a2, b2, c2 float64) {
		a, b, c = f(a, b, c)

		return round(a, precision), round(b, precision), round(c, precision)
	}
}

// SafeFunc is returned by the SafeTo methods of the CoordinateReferenceSystem's
// in this package and the SafeTransform function.
type SafeFunc func(a, b, c float64) (a2, b2, c2 float64, err error)

// Round can round the resulting values to a specific precision.
func (f SafeFunc) Round(precision float64) SafeFunc {
	return func(a, b, c float64) (a2, b2, c2 float64, err error) {
		a, b, c, err = f(a, b, c)

		return round(a, precision), round(b, precision), round(c, precision), err
	}
}
