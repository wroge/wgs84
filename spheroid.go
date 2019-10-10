package wgs84

import "math"

// GRS80 Spheroid.
func GRS80() Spheroid {
	return Spheroid{
		A:  6378137,
		Fi: 298.257222101,
	}
}

// Bessel Spheroid.
func Bessel() Spheroid {
	return Spheroid{
		A:  6377397.155,
		Fi: 299.1528128,
	}
}

// Airy Spheroid.
func Airy() Spheroid {
	return Spheroid{
		A:  6377563.396,
		Fi: 299.3249646,
	}
}

// Hayford Spheroid.
func Hayford() Spheroid {
	return Spheroid{
		A:  6378388,
		Fi: 297,
	}
}

// Krassowsky Spheroid.
func Krassowsky() Spheroid {
	return Spheroid{
		A:  6378245,
		Fi: 298.3,
	}
}

// Spheroid implements the GeodeticSpheroid interface.
// It has helper methods to simplify the implementation of your own
// CoordinateReference system, for example through the Projection struct.
type Spheroid struct {
	A, Fi float64
}

// MajorAxis is a method to implement the GeodeticSpheroid interface.
func (s Spheroid) MajorAxis() float64 {
	if s == (Spheroid{}) {
		return 6378137
	}
	return s.A
}

// InverseFlattening is a method to implement the GeodeticSpheroid interface.
func (s Spheroid) InverseFlattening() float64 {
	if s == (Spheroid{}) {
		return 298.257223563
	}
	return s.Fi
}

// A2 is a helper method.
func (s Spheroid) A2() float64 {
	return s.MajorAxis() * s.MajorAxis()
}

// Fi2 is a helper method.
func (s Spheroid) Fi2() float64 {
	return s.InverseFlattening() * s.InverseFlattening()
}

// F is a helper method.
func (s Spheroid) F() float64 {
	return 1 / s.InverseFlattening()
}

// F2 is a helper method.
func (s Spheroid) F2() float64 {
	return s.F() * s.F()
}

// B is a helper method.
func (s Spheroid) B() float64 {
	return s.MajorAxis() * (1 - s.F())
}

// E2 is a helper method.
func (s Spheroid) E2() float64 {
	return 2/s.InverseFlattening() - s.F2()
}

// E is a helper method.
func (s Spheroid) E() float64 {
	return math.Sqrt(s.E2())
}

// E4 is a helper method.
func (s Spheroid) E4() float64 {
	return s.E2() * s.E2()
}

// E6 is a helper method.
func (s Spheroid) E6() float64 {
	return s.E4() * s.E2()
}

// Ei is a helper method.
func (s Spheroid) Ei() float64 {
	return (1 - math.Sqrt(1-s.E2())) / (1 + math.Sqrt(1-s.E2()))
}

// Ei2 is a helper method.
func (s Spheroid) Ei2() float64 {
	return s.Ei() * s.Ei()
}

// Ei3 is a helper method.
func (s Spheroid) Ei3() float64 {
	return s.Ei2() * s.Ei()
}

// Ei4 is a helper method.
func (s Spheroid) Ei4() float64 {
	return s.Ei3() * s.Ei()
}
