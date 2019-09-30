package system

import (
	"math"

	"github.com/wroge/wgs84/spheroid"
)

type Spheroid = interface {
	A() float64
	Fi() float64
}

func getSpheroid(s Spheroid) spheroid.Spheroid {
	if s == nil {
		return spheroid.WGS84()
	} else if s, ok := s.(spheroid.Spheroid); ok {
		return s
	}
	return spheroid.New(s.A(), s.Fi())
}

func sin2(east float64) float64 {
	return math.Pow(math.Sin(east), 2)
}

func cos2(east float64) float64 {
	return math.Pow(math.Cos(east), 2)
}

func tan2(east float64) float64 {
	return math.Pow(math.Tan(east), 2)
}

func degree(r float64) float64 {
	return r * 180 / math.Pi
}

func radian(d float64) float64 {
	return d * math.Pi / 180
}
