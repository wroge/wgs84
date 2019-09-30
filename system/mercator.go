package system

import (
	"math"
)

type Mercator struct {
	Lonf, Scale, Eastf, Northf float64
}

func (sys Mercator) ToXYZ(east, north, h float64, s Spheroid) (x, y, z float64) {
	sph := getSpheroid(s)
	east = (east - sys.Eastf) / sys.Scale
	north = (north - sys.Northf) / sys.Scale
	t := math.Exp(-north * sph.A())
	φ := math.Pi/2 - 2*math.Atan(t)
	for i := 0; i < 5; i++ {
		φ = math.Pi/2 - 2*math.Atan(t*math.Pow((1-sph.E()*math.Sin(φ))/(1+sph.E()*math.Sin(φ)), sph.E()/2))
	}
	return LonLat{}.ToXYZ(east/sph.A()+sys.Lonf, degree(φ), h, sph)
}

func (sys Mercator) FromXYZ(x, y, z float64, s Spheroid) (east, north, h float64) {
	sph := getSpheroid(s)
	lon, lat, h := LonLat{}.FromXYZ(x, y, z, sph)
	east = sys.Scale * sph.A() * (radian(lon) - radian(sys.Lonf))
	north = sys.Scale * sph.A() / 2 *
		math.Log(1+math.Sin(radian(lat))/(1-math.Sin(radian(lat)))*
			math.Pow((1-sph.E()*math.Sin(radian(lat)))/(1+sph.E()*math.Sin(radian(lat))), math.E))
	return
}
