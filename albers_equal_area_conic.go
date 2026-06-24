package wgs84

import (
	"fmt"
	"math"
)

type AlbersConicEqualArea struct {
	Lonf, Latf, Sp1, Sp2, Eastf, Northf float64
}

func (a AlbersConicEqualArea) String() string {
	return fmt.Sprintf("+proj=aea +lat_0=%s +lon_0=%s +lat_1=%s +lat_2=%s +x_0=%s +y_0=%s",
		projFloat(a.Latf), projFloat(a.Lonf), projFloat(a.Sp1), projFloat(a.Sp2), projFloat(a.Eastf), projFloat(a.Northf),
	)
}

func (a AlbersConicEqualArea) FromGeographic(lon float64, lat float64, h float64, s Spheroid) (east, north, h2 float64) {
	phif := radian(a.Latf)
	phi1 := radian(a.Sp1)
	phi2 := radian(a.Sp2)
	lambdaf := radian(a.Lonf)
	alphaf := (1 - s.E2()) * ((math.Sin(phif) / (1 - s.E2()*sin2(phif))) - (1/(2*s.E()))*math.Log((1-s.E()*math.Sin(phif))/(1+s.E()*math.Sin(phif))))
	alpha1 := (1 - s.E2()) * ((math.Sin(phi1) / (1 - s.E2()*sin2(phi1))) - (1/(2*s.E()))*math.Log((1-s.E()*math.Sin(phi1))/(1+s.E()*math.Sin(phi1))))
	alpha2 := (1 - s.E2()) * ((math.Sin(phi2) / (1 - s.E2()*sin2(phi2))) - (1/(2*s.E()))*math.Log((1-s.E()*math.Sin(phi2))/(1+s.E()*math.Sin(phi2))))

	m1 := math.Cos(phi1) / math.Sqrt(1-s.E2()*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-s.E2()*sin2(phi2))

	n := (m1*m1 - m2*m2) / (alpha2 - alpha1)
	c := m1*m1 + (n * alpha1)
	rf := (s.A * math.Sqrt(c-n*alphaf)) / n

	lambda := radian(lon)
	phi := radian(lat)

	alpha := (1 - s.E2()) * ((math.Sin(phi) / (1 - s.E2()*sin2(phi))) - (1/(2*s.E()))*math.Log((1-s.E()*math.Sin(phi))/(1+s.E()*math.Sin(phi))))

	theta := n * (lambda - lambdaf)
	r := (s.A * math.Sqrt(c-n*alpha)) / n

	east = a.Eastf + r*math.Sin(theta)
	north = a.Northf + rf - r*math.Cos(theta)

	return east, north, h
}

func (a AlbersConicEqualArea) ToGeographic(east, north, h float64, s Spheroid) (lon float64, lat float64, h2 float64) {
	phif := radian(a.Latf)
	phi1 := radian(a.Sp1)
	phi2 := radian(a.Sp2)
	lambdaf := radian(a.Lonf)
	alphaf := (1 - s.E2()) * ((math.Sin(phif) / (1 - s.E2()*sin2(phif))) - (1/(2*s.E()))*math.Log((1-s.E()*math.Sin(phif))/(1+s.E()*math.Sin(phif))))
	alpha1 := (1 - s.E2()) * ((math.Sin(phi1) / (1 - s.E2()*sin2(phi1))) - (1/(2*s.E()))*math.Log((1-s.E()*math.Sin(phi1))/(1+s.E()*math.Sin(phi1))))
	alpha2 := (1 - s.E2()) * ((math.Sin(phi2) / (1 - s.E2()*sin2(phi2))) - (1/(2*s.E()))*math.Log((1-s.E()*math.Sin(phi2))/(1+s.E()*math.Sin(phi2))))

	m1 := math.Cos(phi1) / math.Sqrt(1-s.E2()*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-s.E2()*sin2(phi2))

	n := (m1*m1 - m2*m2) / (alpha2 - alpha1)
	c := m1*m1 + (n * alpha1)
	rf := (s.A * math.Sqrt(c-n*alphaf)) / n

	ri := math.Sqrt(intPow(east-a.Eastf, 2) + intPow(rf-(north-a.Northf), 2))
	alphai := (c - (intPow(ri, 2) * intPow(n, 2) / s.A2())) / n
	betai := math.Asin(alphai / (1 - ((1 - s.E2()) / (2 * s.E()) * math.Log((1-s.E())/(1+s.E())))))

	var theta float64

	if n > 0 {
		theta = math.Atan2((east - a.Eastf), (rf - (north - a.Northf)))
	} else {
		theta = math.Atan2(-(east - a.Eastf), -(rf - (north - a.Northf)))
	}

	phi := betai + ((s.E2()/3 + 31*s.E4()/180 + 517*s.E6()/5040) * math.Sin(2*betai)) + ((23*s.E4()/360 + 251*s.E6()/3780) * math.Sin(4*betai)) + ((761 * s.E6() / 45360) * math.Sin(6*betai))
	lambda := lambdaf + (theta / n)

	return degree(lambda), degree(phi), h
}
