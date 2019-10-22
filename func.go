package wgs84

import "math"

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
