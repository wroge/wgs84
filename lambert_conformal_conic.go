package wgs84

import "math"

type LambertConformalConic2SP struct {
	Lonf, Latf, Sp1, Sp2, Eastf, Northf float64
}

func (l LambertConformalConic2SP) FromGeographic(lon float64, lat float64, h float64, s Spheroid) (east, north, h2 float64) {
	phif := radian(l.Latf)
	phi1 := radian(l.Sp1)
	phi2 := radian(l.Sp2)
	lambdaf := radian(l.Lonf)

	tf := math.Tan(math.Pi/4-phif/2) / math.Pow((1-s.E()*math.Sin(phif))/(1+s.E()*math.Sin(phif)), s.E()/2)
	t1 := math.Tan(math.Pi/4-phi1/2) / math.Pow((1-s.E()*math.Sin(phi1))/(1+s.E()*math.Sin(phi1)), s.E()/2)
	t2 := math.Tan(math.Pi/4-phi2/2) / math.Pow((1-s.E()*math.Sin(phi2))/(1+s.E()*math.Sin(phi2)), s.E()/2)

	m1 := math.Cos(phi1) / math.Sqrt(1-s.E2()*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-s.E2()*sin2(phi2))

	var n float64
	if math.Abs(phi1-phi2) < 1e-14 {
		n = math.Sin(phi1)
	} else {
		n = math.Log(m1/m2) / math.Log(t1/t2)
	}

	f := m1 / (n * math.Pow(t1, n))
	rf := s.A * f * math.Pow(tf, n)

	phi := radian(lat)
	lambda := radian(lon)

	t := math.Tan(math.Pi/4-phi/2) / math.Pow((1-s.E()*math.Sin(phi))/(1+s.E()*math.Sin(phi)), s.E()/2)

	r := s.A * f * math.Pow(t, n)
	theta := n * (lambda - lambdaf)

	east = l.Eastf + r*math.Sin(theta)
	north = l.Northf + rf - r*math.Cos(theta)

	return east, north, h
}

func (l LambertConformalConic2SP) ToGeographic(east, north, h float64, s Spheroid) (lon float64, lat float64, h2 float64) {
	phif := radian(l.Latf)
	phi1 := radian(l.Sp1)
	phi2 := radian(l.Sp2)
	lambdaf := radian(l.Lonf)

	tf := math.Tan(math.Pi/4-phif/2) / math.Pow((1-s.E()*math.Sin(phif))/(1+s.E()*math.Sin(phif)), s.E()/2)
	t1 := math.Tan(math.Pi/4-phi1/2) / math.Pow((1-s.E()*math.Sin(phi1))/(1+s.E()*math.Sin(phi1)), s.E()/2)
	t2 := math.Tan(math.Pi/4-phi2/2) / math.Pow((1-s.E()*math.Sin(phi2))/(1+s.E()*math.Sin(phi2)), s.E()/2)

	m1 := math.Cos(phi1) / math.Sqrt(1-s.E2()*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-s.E2()*sin2(phi2))

	var n float64
	if math.Abs(phi1-phi2) < 1e-14 {
		n = math.Sin(phi1)
	} else {
		n = math.Log(m1/m2) / math.Log(t1/t2)
	}

	f := m1 / (n * math.Pow(t1, n))
	rf := s.A * f * math.Pow(tf, n)

	ri := math.Hypot(east-l.Eastf, rf-(north-l.Northf))

	ti := math.Pow(ri/(s.A*f), 1/n)

	theta := math.Atan2(east-l.Eastf, rf-(north-l.Northf))

	phi := math.Pi/2 - 2*math.Atan(ti)

	for range 6 {
		next := math.Pi/2 - 2*math.Atan(ti*math.Pow((1-s.E()*math.Sin(phi))/(1+s.E()*math.Sin(phi)), s.E()/2))
		if math.Abs(next-phi) < 1e-12 {
			phi = next
			break
		}

		phi = next
	}

	lambda := theta/n + lambdaf

	return degree(lambda), degree(phi), h
}

type LambertConformalConic1SP struct {
	Lonf, Latf, Scale, Eastf, Northf float64
}

func (l LambertConformalConic1SP) FromGeographic(lon float64, lat float64, h float64, s Spheroid) (east, north, h2 float64) {
	phif := radian(l.Latf)
	lambdaf := radian(l.Lonf)

	tf := math.Tan(math.Pi/4-phif/2) / math.Pow((1-s.E()*math.Sin(phif))/(1+s.E()*math.Sin(phif)), s.E()/2)
	m1 := math.Cos(phif) / math.Sqrt(1-s.E2()*sin2(phif))

	n := math.Sin(phif)
	f := l.Scale * m1 / (n * math.Pow(tf, n))
	rf := s.A * f * math.Pow(tf, n)

	phi := radian(lat)
	lambda := radian(lon)

	t := math.Tan(math.Pi/4-phi/2) / math.Pow((1-s.E()*math.Sin(phi))/(1+s.E()*math.Sin(phi)), s.E()/2)

	r := s.A * f * math.Pow(t, n)
	theta := n * (lambda - lambdaf)

	east = l.Eastf + r*math.Sin(theta)
	north = l.Northf + rf - r*math.Cos(theta)

	return east, north, h
}

func (l LambertConformalConic1SP) ToGeographic(east, north, h float64, s Spheroid) (lon float64, lat float64, h2 float64) {
	phif := radian(l.Latf)
	lambdaf := radian(l.Lonf)

	tf := math.Tan(math.Pi/4-phif/2) / math.Pow((1-s.E()*math.Sin(phif))/(1+s.E()*math.Sin(phif)), s.E()/2)
	m1 := math.Cos(phif) / math.Sqrt(1-s.E2()*sin2(phif))

	n := math.Sin(phif)
	f := l.Scale * m1 / (n * math.Pow(tf, n))
	rf := s.A * f * math.Pow(tf, n)
	ri := math.Hypot(east-l.Eastf, rf-(north-l.Northf))

	ti := math.Pow(ri/(s.A*f), 1/n)

	theta := math.Atan2(east-l.Eastf, rf-(north-l.Northf))

	phi := math.Pi/2 - 2*math.Atan(ti)

	for range 6 {
		next := math.Pi/2 - 2*math.Atan(ti*math.Pow((1-s.E()*math.Sin(phi))/(1+s.E()*math.Sin(phi)), s.E()/2))
		if math.Abs(next-phi) < 1e-12 {
			phi = next
			break
		}

		phi = next
	}

	lambda := theta/n + lambdaf

	return degree(lambda), degree(phi), h
}
