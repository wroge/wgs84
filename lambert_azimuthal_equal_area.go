package wgs84

import (
	"fmt"
	"math"
)

type LambertAzimuthalEqualArea struct {
	Lonf, Latf, Eastf, Northf float64
}

func (l LambertAzimuthalEqualArea) String() string {
	return fmt.Sprintf("+proj=laea +lat_0=%g +lon_0=%g +x_0=%g +y_0=%g",
		l.Latf, l.Lonf, l.Eastf, l.Northf,
	)
}

func (l LambertAzimuthalEqualArea) FromGeographic(lon float64, lat float64, h float64, s Spheroid) (east, north, h2 float64) {
	phi0 := radian(l.Latf)
	lambda0 := radian(l.Lonf)

	q0 := (1 - s.E2()) * ((math.Sin(phi0) / (1 - s.E2()*sin2(phi0))) - (1 / (2 * s.E()) * math.Log((1-s.E()*math.Sin(phi0))/(1+s.E()*math.Sin(phi0)))))
	qp := (1 - s.E2()) * ((1 / (1 - s.E2())) - ((1 / (2 * s.E())) * math.Log((1-s.E())/(1+s.E()))))

	beta0 := math.Asin(q0 / qp)
	rq := s.A * math.Sqrt(qp/2)
	g := s.A * (math.Cos(phi0) / math.Sqrt(1-s.E2()*sin2(phi0))) / (rq * math.Cos(beta0))

	phi := radian(lat)
	lambda := radian(lon)

	q := (1 - s.E2()) * ((math.Sin(phi) / (1 - s.E2()*sin2(phi))) - (1 / (2 * s.E()) * math.Log((1-s.E()*math.Sin(phi))/(1+s.E()*math.Sin(phi)))))

	beta := math.Asin(q / qp)
	b := rq * math.Sqrt(2/(1+math.Sin(beta0)*math.Sin(beta)+(math.Cos(beta0)*math.Cos(beta)*math.Cos(lambda-lambda0))))

	east = l.Eastf + ((b * g) * (math.Cos(beta) * math.Sin(lambda-lambda0)))
	north = l.Northf + (b/g)*((math.Cos(beta0)*math.Sin(beta))-(math.Sin(beta0)*math.Cos(beta)*math.Cos(lambda-lambda0)))

	return east, north, h
}

func (l LambertAzimuthalEqualArea) ToGeographic(east, north, h float64, s Spheroid) (lon float64, lat float64, h2 float64) {
	phi0 := radian(l.Latf)
	lambda0 := radian(l.Lonf)

	q0 := (1 - s.E2()) * ((math.Sin(phi0) / (1 - s.E2()*sin2(phi0))) - (1 / (2 * s.E()) * math.Log((1-s.E()*math.Sin(phi0))/(1+s.E()*math.Sin(phi0)))))
	qp := (1 - s.E2()) * ((1 / (1 - s.E2())) - ((1 / (2 * s.E())) * math.Log((1-s.E())/(1+s.E()))))

	beta0 := math.Asin(q0 / qp)
	rq := s.A * math.Sqrt(qp/2)
	g := s.A * (math.Cos(phi0) / math.Sqrt(1-s.E2()*sin2(phi0))) / (rq * math.Cos(beta0))

	rho := math.Sqrt(intPow((east-l.Eastf)/g, 2) + intPow(g*(north-l.Northf), 2))
	c := 2 * math.Asin(rho/(2*rq))
	betai := math.Asin((math.Cos(c) * math.Sin(beta0)) + ((g * (north - l.Northf) * math.Sin(c) * math.Cos(beta0)) / rho))

	phi := betai + ((s.E2()/3 + (31*s.E4())/180 + (517*s.E6())/5040) * math.Sin(2*betai)) +
		((23*s.E4())/360+(251*s.E6())/3780)*math.Sin(4*betai) +
		((761*s.E6())/45360)*math.Sin(6*betai)
	lambda := lambda0 + math.Atan2((east-l.Eastf)*math.Sin(c), (g*rho*math.Cos(beta0)*math.Cos(c)-intPow(g, 2)*(north-l.Northf)*math.Sin(beta0)*math.Sin(c)))

	return degree(lambda), degree(phi), h
}
