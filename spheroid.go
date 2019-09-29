package wgs84

import "math"

type spheroid struct {
	a, fi float64
}

func (s spheroid) A() float64 {
	return s.a
}

func (s spheroid) Fi() float64 {
	return s.fi
}

func (s spheroid) A2() float64 {
	return s.A() * s.A()
}

func (s spheroid) Fi2() float64 {
	return s.Fi() * s.Fi()
}

func (s spheroid) F() float64 {
	return 1 / s.Fi()
}

func (s spheroid) F2() float64 {
	return s.F() * s.F()
}

func (s spheroid) B() float64 {
	return s.A() * (1 - s.F())
}

func (s spheroid) E2() float64 {
	return 2/s.Fi() - s.F2()
}

func (s spheroid) E() float64 {
	return math.Sqrt(s.E2())
}

func (s spheroid) E4() float64 {
	return s.E2() * s.E2()
}

func (s spheroid) E6() float64 {
	return s.E4() * s.E2()
}

func (s spheroid) Ei() float64 {
	return (1 - math.Sqrt(1-s.E2())) / (1 + math.Sqrt(1-s.E2()))
}

func (s spheroid) Ei2() float64 {
	return s.Ei() * s.Ei()
}

func (s spheroid) Ei3() float64 {
	return s.Ei2() * s.Ei()
}

func (s spheroid) Ei4() float64 {
	return s.Ei3() * s.Ei()
}
