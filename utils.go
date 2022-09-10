//nolint:nonamedreturns,varnamelen,gomnd,asciicheck,nosnakecase
package wgs84

import "math"

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

func lonLatToXYZ(lon, lat, h, a, fi float64) (x, y, z float64) {
	s := spheroid{a: a, fi: fi}
	x = (_N(radian(lat), s) + h) * math.Cos(radian(lon)) * math.Cos(radian(lat))
	y = (_N(radian(lat), s) + h) * math.Cos(radian(lat)) * math.Sin(radian(lon))
	z = (_N(radian(lat), s)*math.Pow(s.A()*(1-s.f()), 2)/(s.a2()) + h) * math.Sin(radian(lat))

	return x, y, z
}

func xyzToLonLat(x, y, z, a, fi float64) (lon, lat, h float64) {
	s := spheroid{a: a, fi: fi}
	sd := math.Sqrt(x*x + y*y)
	T := math.Atan(z * s.A() / (sd * s.b()))
	B := math.Atan((z + s.e2()*(s.a2())/s.b()*
		math.Pow(math.Sin(T), 3)) / (sd - s.e2()*s.A()*math.Pow(math.Cos(T), 3)))
	h = sd/math.Cos(B) - _N(B, s)
	lon = degree(math.Atan2(y, x))
	lat = degree(B)

	return lon, lat, h
}

func _N(φ float64, s spheroid) float64 {
	return s.A() / math.Sqrt(1-s.e2()*math.Pow(math.Sin(φ), 2))
}
