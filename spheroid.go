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

func (s spheroid) a2() float64 {
	return s.A() * s.A()
}

func (s spheroid) f() float64 {
	return 1 / s.Fi()
}

func (s spheroid) f2() float64 {
	return s.f() * s.f()
}

func (s spheroid) b() float64 {
	return s.A() * (1 - s.f())
}

func (s spheroid) e2() float64 {
	return 2/s.Fi() - s.f2()
}

func (s spheroid) e() float64 {
	return math.Sqrt(s.e2())
}

func (s spheroid) e4() float64 {
	return s.e2() * s.e2()
}

func (s spheroid) e6() float64 {
	return s.e4() * s.e2()
}

func (s spheroid) ei() float64 {
	return (1 - math.Sqrt(1-s.e2())) / (1 + math.Sqrt(1-s.e2()))
}

func (s spheroid) ei2() float64 {
	return s.ei() * s.ei()
}

func (s spheroid) ei3() float64 {
	return s.ei2() * s.ei()
}

func (s spheroid) ei4() float64 {
	return s.ei3() * s.ei()
}
