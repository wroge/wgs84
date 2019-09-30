package spheroid

import "math"

func WGS84() Spheroid {
	return Spheroid{6378137, 298.257223563}
}

func GRS80() Spheroid {
	return Spheroid{6378137, 298.257222101}
}

func Bessel() Spheroid {
	return Spheroid{6377397.155, 299.1528128}
}

func Airy() Spheroid {
	return Spheroid{6377563.396, 299.3249646}
}

func Hayford() Spheroid {
	return Spheroid{6378388, 297}
}

func Krassowsky() Spheroid {
	return Spheroid{6378245, 298.3}
}

func New(a, fi float64) Spheroid {
	return Spheroid{a, fi}
}

type Spheroid struct {
	a, fi float64
}

func (s Spheroid) A() float64 {
	if s.a == 0 {
		return WGS84().A()
	}
	return s.a
}

func (s Spheroid) Fi() float64 {
	if s.fi == 0 {
		return WGS84().Fi()
	}
	return s.fi
}

func (s Spheroid) A2() float64 {
	return s.A() * s.A()
}

func (s Spheroid) Fi2() float64 {
	return s.Fi() * s.Fi()
}

func (s Spheroid) F() float64 {
	return 1 / s.Fi()
}

func (s Spheroid) F2() float64 {
	return s.F() * s.F()
}

func (s Spheroid) B() float64 {
	return s.A() * (1 - s.F())
}

func (s Spheroid) E2() float64 {
	return 2/s.Fi() - s.F2()
}

func (s Spheroid) E() float64 {
	return math.Sqrt(s.E2())
}

func (s Spheroid) E4() float64 {
	return s.E2() * s.E2()
}

func (s Spheroid) E6() float64 {
	return s.E4() * s.E2()
}

func (s Spheroid) Ei() float64 {
	return (1 - math.Sqrt(1-s.E2())) / (1 + math.Sqrt(1-s.E2()))
}

func (s Spheroid) Ei2() float64 {
	return s.Ei() * s.Ei()
}

func (s Spheroid) Ei3() float64 {
	return s.Ei2() * s.Ei()
}

func (s Spheroid) Ei4() float64 {
	return s.Ei3() * s.Ei()
}
