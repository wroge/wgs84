package wgs84

import (
	"fmt"
	"math"
)

type SwissObliqueMercator struct {
	Lonf, Latf, Scale, Eastf, Northf float64
}

func (som SwissObliqueMercator) String() string {
	return fmt.Sprintf("+proj=somerc +lat_0=%s +lon_0=%s +k_0=%s +x_0=%s +y_0=%s",
		projFloat(som.Latf), projFloat(som.Lonf), projFloat(som.Scale), projFloat(som.Eastf), projFloat(som.Northf),
	)
}

func (som SwissObliqueMercator) FromGeographic(lon float64, lat float64, height float64, s Spheroid) (east, north, height2 float64) {
	a := s.A
	e := s.E()
	e2 := s.E2()

	phiC := radian(som.Latf)
	lambdaC := radian(som.Lonf)
	alphaC := radian(90)
	gammaC := radian(90)

	sinφ := math.Sin(phiC)
	cosφ := math.Cos(phiC)

	sLat := sign(sinφ)

	B := math.Sqrt(1 + (e2*math.Pow(cosφ, 4))/(1-e2))

	A := a * som.Scale * math.Sqrt(1-e2) * B /
		(1 - e2*sinφ*sinφ)

	t0 := math.Tan(math.Pi/4-phiC/2) /
		math.Pow((1-e*sinφ)/(1+e*sinφ), e/2)

	D := (B * math.Sqrt(1-e2)) /
		(cosφ * math.Sqrt(1-e2*sinφ*sinφ))

	D2 := D * D
	if D < 1 {
		D2 = 1
	}

	sqrtTerm := math.Sqrt(D2 - 1)

	F := D + sqrtTerm*sLat
	G := (F - 1/F) / 2

	x := math.Sin(alphaC) / D
	x = clamp(x, -1, 1)

	gamma0 := math.Asin(x)

	lonArg := clamp(G*math.Tan(gamma0), -1, 1)
	lon0 := lambdaC - math.Asin(lonArg)/B

	H := F * math.Pow(t0, B)

	var uc float64
	if math.Abs(math.Abs(alphaC)-math.Pi/2) < 1e-12 {
		uc = A * (lambdaC - lon0)
	} else {
		uc = (A / B) * math.Atan2(sqrtTerm, math.Cos(alphaC)) * sLat
	}

	phi := radian(lat)
	lambda := radian(lon)

	t := math.Tan(math.Pi/4-phi/2) /
		math.Pow((1-e*math.Sin(phi))/(1+e*math.Sin(phi)), e/2)

	Q := H / math.Pow(t, B)

	S := (Q - 1/Q) / 2
	T := (Q + 1/Q) / 2

	V := math.Sin(B * (lambda - lon0))
	U := (-V*math.Cos(gamma0) + S*math.Sin(gamma0)) / T

	v := (A / (2 * B)) * math.Log((1-U)/(1+U))

	u0 := (A / B) * math.Atan2(
		S*math.Cos(gamma0)+V*math.Sin(gamma0),
		math.Cos(B*(lambda-lon0)),
	)

	u := u0 - math.Abs(uc)*sLat*sign(lambda-lon0)

	E := v*math.Cos(gammaC) + u*math.Sin(gammaC) + som.Eastf
	N := u*math.Cos(gammaC) - v*math.Sin(gammaC) + som.Northf

	return E, N, height
}

func (som SwissObliqueMercator) ToGeographic(east, north, height float64, s Spheroid) (lon float64, lat float64, height2 float64) {
	a := s.A
	e := s.E()
	e2 := s.E2()

	phiC := radian(som.Latf)
	lambdaC := radian(som.Lonf)
	alphaC := radian(90)
	gammaC := radian(90)

	sinφ := math.Sin(phiC)
	cosφ := math.Cos(phiC)

	sLat := sign(sinφ)

	B := math.Sqrt(1 + (e2*math.Pow(cosφ, 4))/(1-e2))

	A := a * som.Scale * math.Sqrt(1-e2) * B /
		(1 - e2*sinφ*sinφ)

	t0 := math.Tan(math.Pi/4-phiC/2) /
		math.Pow((1-e*sinφ)/(1+e*sinφ), e/2)

	D := (B * math.Sqrt(1-e2)) /
		(cosφ * math.Sqrt(1-e2*sinφ*sinφ))

	D2 := D * D
	if D < 1 {
		D2 = 1
	}

	sqrtTerm := math.Sqrt(D2 - 1)

	F := D + sqrtTerm*sLat
	G := (F - 1/F) / 2

	x := math.Sin(alphaC) / D
	x = clamp(x, -1, 1)

	gamma0 := math.Asin(x)

	lonArg := clamp(G*math.Tan(gamma0), -1, 1)
	lon0 := lambdaC - math.Asin(lonArg)/B

	H := F * math.Pow(t0, B)

	var uc float64
	if math.Abs(math.Abs(alphaC)-math.Pi/2) < 1e-12 {
		uc = A * (lambdaC - lon0)
	} else {
		uc = (A / B) * math.Atan2(sqrtTerm, math.Cos(alphaC)) * sLat
	}

	v := (east-som.Eastf)*math.Cos(gammaC) -
		(north-som.Northf)*math.Sin(gammaC)

	u := (north-som.Northf)*math.Cos(gammaC) +
		(east-som.Eastf)*math.Sin(gammaC)

	u += math.Abs(uc) * sLat

	Qp := math.Exp(-(B * v / A))

	Sp := 0.5 * (Qp - 1/Qp)
	Tp := 0.5 * (Qp + 1/Qp)

	Vp := math.Sin(B * u / A)

	Up := (Vp*math.Cos(gamma0) + Sp*math.Sin(gamma0)) / Tp

	tp := math.Pow(H/math.Sqrt((1+Up)/(1-Up)), 1.0/B)

	chi := math.Pi/2 - 2*math.Atan(tp)

	lat = chi +
		math.Sin(2*chi)*(e*e/2+5*math.Pow(e, 4)/24+math.Pow(e, 6)/12+13*math.Pow(e, 8)/360) +
		math.Sin(4*chi)*(7*math.Pow(e, 4)/48+29*math.Pow(e, 6)/240+811*math.Pow(e, 8)/11520) +
		math.Sin(6*chi)*(7*math.Pow(e, 6)/120+81*math.Pow(e, 8)/1120) +
		math.Sin(8*chi)*(4279*math.Pow(e, 8)/161280)

	lon = lon0 -
		(1.0/B)*math.Atan2(
			Sp*math.Cos(gamma0)-Vp*math.Sin(gamma0),
			math.Cos(B*u/A),
		)

	return degree(lon), degree(lat), 0
}
