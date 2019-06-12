// Package spheroid implements the Spheroid interface
// and some useful methods for coordinate conversions.
package spheroid

import "math"

// WGS84 Spheroid used mostly by the WGS84 geodetic datum.
func WGS84() Spheroid {
	return Spheroid{6378137, 298.257223563}
}

// GRS80 Spheroid is used by ETRS89, among others.
func GRS80() Spheroid {
	return Spheroid{6378137, 298.257222101}
}

// Bessel Spheroid is used by DHDN2001, among others.
func Bessel() Spheroid {
	return Spheroid{6377397.155, 299.1528128}
}

// Airy Spheroid is used by OSGB36, among others.
func Airy() Spheroid {
	return Spheroid{6377563.396, 299.3249646}
}

// Hayford Spheroid is used by many geodetic datums.
func Hayford() Spheroid {
	return Spheroid{6378388, 297}
}

// Krassowsky Spheroid was mostly used in the former
// Soviet republics.
func Krassowsky() Spheroid {
	return Spheroid{6378245, 298.3}
}

// New is a constructor for Spheroids.
func New(a, fi float64) Spheroid {
	return Spheroid{a, fi}
}

// Spheroid has to be constructed by New.
type Spheroid struct {
	a, fi float64
}

// A is used by the Spheroid interface used in wgs84.System.
func (s Spheroid) A() float64 {
	if s.a == 0 {
		return WGS84().A()
	}
	return s.a
}

// Fi is used by the Spheroid interface used in wgs84.System.
func (s Spheroid) Fi() float64 {
	if s.fi == 0 {
		return WGS84().Fi()
	}
	return s.fi
}

// A2 is a helper function for coordinate conversion.
func (s Spheroid) A2() float64 {
	return s.A() * s.A()
}

// Fi2 is a helper function for coordinate conversion.
func (s Spheroid) Fi2() float64 {
	return s.Fi() * s.Fi()
}

// F is a helper function for coordinate conversion.
func (s Spheroid) F() float64 {
	return 1 / s.Fi()
}

// F2 is a helper function for coordinate conversion.
func (s Spheroid) F2() float64 {
	return s.F() * s.F()
}

// B is a helper function for coordinate conversion.
func (s Spheroid) B() float64 {
	return s.A() * (1 - s.F())
}

// E2 is a helper function for coordinate conversion.
func (s Spheroid) E2() float64 {
	return 2/s.Fi() - s.F2()
}

// E is a helper function for coordinate conversion.
func (s Spheroid) E() float64 {
	return math.Sqrt(s.E2())
}

// E4 is a helper function for coordinate conversion.
func (s Spheroid) E4() float64 {
	return s.E2() * s.E2()
}

// E6 is a helper function for coordinate conversion.
func (s Spheroid) E6() float64 {
	return s.E4() * s.E2()
}

// Ei is a helper function for coordinate conversion.
func (s Spheroid) Ei() float64 {
	return (1 - math.Sqrt(1-s.E2())) / (1 + math.Sqrt(1-s.E2()))
}

// Ei2 is a helper function for coordinate conversion.
func (s Spheroid) Ei2() float64 {
	return s.Ei() * s.Ei()
}

// Ei3 is a helper function for coordinate conversion.
func (s Spheroid) Ei3() float64 {
	return s.Ei2() * s.Ei()
}

// Ei4 is a helper function for coordinate conversion.
func (s Spheroid) Ei4() float64 {
	return s.Ei3() * s.Ei()
}
