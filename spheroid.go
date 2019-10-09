package wgs84

import "math"

func GRS80() Spheroid {
	return Spheroid{
		A:  6378137,
		Fi: 298.257222101,
	}
}

func Bessel() Spheroid {
	return Spheroid{
		A:  6377397.155,
		Fi: 299.1528128,
	}
}

func Airy() Spheroid {
	return Spheroid{
		A:  6377563.396,
		Fi: 299.3249646,
	}
}

func Hayford() Spheroid {
	return Spheroid{
		A:  6378388,
		Fi: 297,
	}
}

func Krassowsky() Spheroid {
	return Spheroid{
		A:  6378245,
		Fi: 298.3,
	}
}

type Spheroid struct {
	A, Fi float64
}

func (s Spheroid) MajorAxis() float64 {
	if s == (Spheroid{}) {
		return 6378137
	}
	return s.A
}

func (s Spheroid) InverseFlattening() float64 {
	if s == (Spheroid{}) {
		return 298.257223563
	}
	return s.Fi
}

func (s Spheroid) A2() float64 {
	return s.MajorAxis() * s.MajorAxis()
}

func (s Spheroid) Fi2() float64 {
	return s.InverseFlattening() * s.InverseFlattening()
}

func (s Spheroid) F() float64 {
	return 1 / s.InverseFlattening()
}

func (s Spheroid) F2() float64 {
	return s.F() * s.F()
}

func (s Spheroid) B() float64 {
	return s.MajorAxis() * (1 - s.F())
}

func (s Spheroid) E2() float64 {
	return 2/s.InverseFlattening() - s.F2()
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
