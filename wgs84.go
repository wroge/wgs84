package wgs84

import (
	"math"
)

var (
	WGS84      = NewSpheroid(6378137, 298.257223563)
	GRS80      = NewSpheroid(6378137, 298.257222101)
	Airy1830   = NewSpheroid(6377563.396, 299.3249646)
	Bessel1841 = NewSpheroid(6377397.155, 299.1528128)
	Clarke1866 = NewSpheroid(6378206.4, 294.9786982139006)
)

type CoordinateReferenceSystem interface {
	ToWGS84(x, y, z float64) (x0, y0, z0 float64)
	FromWGS84(x0, y0, z0 float64) (x, y, z float64)
}

func Transform(from, to CoordinateReferenceSystem) TransformFunc {
	return func(a, b, c float64) (a2 float64, b2 float64, c2 float64) {
		if from != nil {
			a, b, c = from.ToWGS84(a, b, c)
		}

		if to == nil {
			return a, b, c
		}

		return to.FromWGS84(a, b, c)
	}
}

type TransformFunc func(a, b, c float64) (a2, b2, c2 float64)

type Projection interface {
	ToGeographic(east, north float64, s Spheroid) (lon, lat float64)
	FromGeographic(lon, lat float64, s Spheroid) (east, north float64)
}

type Geographic struct {
	Geocentric CoordinateReferenceSystem
	Spheroid   Spheroid
}

func (g Geographic) ToWGS84(lon, lat, h float64) (x0, y0, z0 float64) {
	x, y, z := g.Spheroid.ToGeocentric(lon, lat, h)

	if g.Geocentric == nil {
		return x, y, z
	}

	return g.Geocentric.ToWGS84(x, y, z)
}

func (g Geographic) FromWGS84(x0, y0, z0 float64) (a, b, c float64) {
	if g.Geocentric != nil {
		x0, y0, z0 = g.Geocentric.FromWGS84(x0, y0, z0)
	}

	return g.Spheroid.FromGeocentric(x0, y0, z0)
}

type Projected struct {
	Geographic Geographic
	Projection Projection
}

func (p Projected) ToWGS84(east, north, h float64) (x0, y0, z0 float64) {
	lon, lat := p.Projection.ToGeographic(east, north, p.Geographic.Spheroid)

	return p.Geographic.ToWGS84(lon, lat, h)
}

func (p Projected) FromWGS84(x0, y0, z0 float64) (east, north, h float64) {
	lon, lat, h := p.Geographic.FromWGS84(x0, y0, z0)

	east, north = p.Projection.FromGeographic(lon, lat, p.Geographic.Spheroid)

	return east, north, h
}

type root struct{}

func (root) ToWGS84(a, b, c float64) (x0, y0, z0 float64) {
	return a, b, c
}

func (root) FromWGS84(x0, y0, z0 float64) (a, b, c float64) {
	return x0, y0, z0
}

type Helmert struct {
	Tx, Ty, Tz, Rx, Ry, Rz, Ds float64
}

func (t Helmert) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return calcHelmert(x, y, z, t.Tx, t.Ty, t.Tz, t.Rx, t.Ry, t.Rz, t.Ds)
}

func (t Helmert) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return calcHelmert(x0, y0, z0, -t.Tx, -t.Ty, -t.Tz, -t.Rx, -t.Ry, -t.Rz, -t.Ds)
}

const (
	asec = math.Pi / 648000
	ppm  = 0.000001
)

func calcHelmert(x, y, z, Tx, Ty, Tz, Rx, Ry, Rz, Ds float64) (x0, y0, z0 float64) {
	x0 = (1+Ds*ppm)*(x+z*Ry*asec-y*Rz*asec) + Tx
	y0 = (1+Ds*ppm)*(y+x*Rz*asec-z*Rx*asec) + Ty
	z0 = (1+Ds*ppm)*(z+y*Rx*asec-x*Ry*asec) + Tz

	return
}

type WebMercator struct{}

func (p WebMercator) ToGeographic(east, north float64, s Spheroid) (lon, lat float64) {
	D := (-north) / s.A
	phi := math.Pi/2 - 2*math.Atan(math.Pow(math.E, D))
	lambda := east / s.A

	return degree(lambda), degree(phi)
}

func (p WebMercator) FromGeographic(lon, lat float64, s Spheroid) (east, north float64) {
	lambda := radian(lon)
	phi := radian(lat)

	east = s.A * lambda
	north = s.A * math.Log(math.Tan(math.Pi/4+phi/2))

	return east, north
}

type TransverseMercator struct {
	CentralMeridian  float64
	LatitudeOfOrigin float64
	ScaleFactor      float64
	FalseEasting     float64
	FalseNorthing    float64
}

func (p TransverseMercator) ToGeographic(east, north float64, s Spheroid) (lon, lat float64) {
	phi0 := radian(p.LatitudeOfOrigin)
	lambda0 := radian(p.CentralMeridian)

	n := s.F / (2 - s.F)
	n2 := math.Pow(n, 2)
	n3 := math.Pow(n, 3)
	n4 := math.Pow(n, 4)
	B := (s.A / (1 + n)) * (1 + n2/4 + n4/64)
	h1 := n/2.0 - (2/3.0)*n2 + (5/16.0)*n3 + (41/180.0)*n4
	h2 := (13/48.0)*n2 - (3/5.0)*n3 + (557/1440.0)*n4
	h3 := (61/240.0)*n3 - (103/140.0)*n4
	h4 := (49561 / 161280.0) * n4

	var M0 float64

	if phi0 == 0 {
		M0 = 0
	} else if phi0 == math.Pi/2 {
		M0 = B * (math.Pi / 2)
	} else if phi0 == -math.Pi/2 {
		M0 = B * (-math.Pi / 2)
	} else {
		Q0 := math.Asinh(math.Tan(phi0)) - (s.E * math.Atanh(s.E*math.Sin(phi0)))
		xi00 := math.Atan(math.Sinh(Q0))
		xi01 := h1 * math.Sin(2*xi00)
		xi02 := h2 * math.Sin(4*xi00)
		xi03 := h3 * math.Sin(6*xi00)
		xi04 := h4 * math.Sin(8*xi00)
		xi0 := xi00 + xi01 + xi02 + xi03 + xi04
		M0 = B * xi0
	}

	h1i := n/2.0 - (2/3.0)*n2 + (37/96.0)*n3 + (1/360.0)*n4
	h2i := (1/48.0)*n2 - (1/15.0)*n3 + (437/1440.0)*n4
	h3i := (17/480.0)*n3 - (37/840.0)*n4
	h4i := (4397 / 161280.0) * n4

	etai := (east - p.FalseEasting) / (B * p.ScaleFactor)
	xii := ((north - p.FalseNorthing) + p.ScaleFactor*M0) / (B * p.ScaleFactor)

	xi1i := h1i * math.Sin(2*xii) * math.Cosh(2*etai)
	xi2i := h2i * math.Sin(4*xii) * math.Cosh(4*etai)
	xi3i := h3i * math.Sin(6*xii) * math.Cosh(6*etai)
	xi4i := h4i * math.Sin(8*xii) * math.Cosh(8*etai)

	xi0i := xii - (xi1i + xi2i + xi3i + xi4i)

	eta1i := h1i * math.Cos(2*xii) * math.Sinh(2*etai)
	eta2i := h2i * math.Cos(4*xii) * math.Sinh(4*etai)
	eta3i := h3i * math.Cos(6*xii) * math.Sinh(6*etai)
	eta4i := h4i * math.Cos(8*xii) * math.Sinh(8*etai)

	eta0i := etai - (eta1i + eta2i + eta3i + eta4i)

	betai := math.Asin(math.Sin(xi0i) / math.Cosh(eta0i))
	Qi := math.Asinh(math.Tan(betai))
	Qii := Qi + (s.E * math.Atanh(s.E*math.Tanh(Qi)))

	for i := 0; i < 15; i++ {
		q := Qi + (s.E * math.Atanh(s.E*math.Tanh(Qii)))

		Qii = q
	}

	phi := math.Atan(math.Sinh(Qii))
	lambda := lambda0 + math.Asin(math.Tanh(eta0i)/math.Cos(betai))

	return degree(lambda), degree(phi)
}

func (p TransverseMercator) FromGeographic(lon, lat float64, s Spheroid) (east, north float64) {
	phi := radian(lat)
	phi0 := radian(p.LatitudeOfOrigin)
	lambda := radian(lon)
	lambda0 := radian(p.CentralMeridian)

	n := s.F / (2 - s.F)
	n2 := math.Pow(n, 2)
	n3 := math.Pow(n, 3)
	n4 := math.Pow(n, 4)
	B := (s.A / (1 + n)) * (1 + n2/4 + n4/64)
	h1 := n/2 - (2/3.0)*n2 + (5/16.0)*n3 + (41/180.0)*n4
	h2 := (13/48.0)*n2 - (3/5.0)*n3 + (557/1440.0)*n4
	h3 := (61/240.0)*n3 - (103/140.0)*n4
	h4 := (49561 / 161280.0) * n4

	var M0 float64

	if phi0 == 0 {
		M0 = 0
	} else if phi0 == math.Pi/2 {
		M0 = B * (math.Pi / 2)
	} else if phi0 == -math.Pi/2 {
		M0 = B * (-math.Pi / 2)
	} else {
		Q0 := math.Asinh(math.Tan(phi0)) - (s.E * math.Atanh(s.E*math.Sin(phi0)))
		xi00 := math.Atan(math.Sinh(Q0))
		xi01 := h1 * math.Sin(2*xi00)
		xi02 := h2 * math.Sin(4*xi00)
		xi03 := h3 * math.Sin(6*xi00)
		xi04 := h4 * math.Sin(8*xi00)
		xi0 := xi00 + xi01 + xi02 + xi03 + xi04
		M0 = B * xi0
	}

	Q := math.Asinh(math.Tan(phi)) - s.E*math.Atanh(s.E*math.Sin(phi))
	beta := math.Atan(math.Sinh(Q))
	eta0 := math.Atanh(math.Cos(beta) * math.Sin(lambda-lambda0))
	xi0 := math.Asin(math.Sin(beta) * math.Cosh(eta0))

	xi1 := h1 * math.Sin(2*xi0) * math.Cosh(2*eta0)
	xi2 := h2 * math.Sin(4*xi0) * math.Cosh(4*eta0)
	xi3 := h3 * math.Sin(6*xi0) * math.Cosh(6*eta0)
	xi4 := h4 * math.Sin(8*xi0) * math.Cosh(8*eta0)

	xi := xi0 + xi1 + xi2 + xi3 + xi4

	eta1 := h1 * math.Cos(2*xi0) * math.Sinh(2*eta0)
	eta2 := h2 * math.Cos(4*xi0) * math.Sinh(4*eta0)
	eta3 := h3 * math.Cos(6*xi0) * math.Sinh(6*eta0)
	eta4 := h4 * math.Cos(8*xi0) * math.Sinh(8*eta0)

	eta := eta0 + eta1 + eta2 + eta3 + eta4

	east = p.FalseEasting + p.ScaleFactor*B*eta
	north = p.FalseNorthing + p.ScaleFactor*(B*xi-M0)

	return east, north
}

type LambertConformalConic2SP struct {
	CentralMeridian   float64
	LatitudeOfOrigin  float64
	StandardParallel1 float64
	StandardParallel2 float64
	FalseEasting      float64
	FalseNorthing     float64
}

func (p LambertConformalConic2SP) ToGeographic(east, north float64, s Spheroid) (lon, lat float64) {
	phif := radian(p.LatitudeOfOrigin)
	phi1 := radian(p.StandardParallel1)
	phi2 := radian(p.StandardParallel2)
	lambdaf := radian(p.CentralMeridian)

	tf := math.Tan(math.Pi/4-phif/2) / math.Pow((1-s.E*math.Sin(phif))/(1+s.E*math.Sin(phif)), s.E/2)
	t1 := math.Tan(math.Pi/4-phi1/2) / math.Pow((1-s.E*math.Sin(phi1))/(1+s.E*math.Sin(phi1)), s.E/2)
	t2 := math.Tan(math.Pi/4-phi2/2) / math.Pow((1-s.E*math.Sin(phi2))/(1+s.E*math.Sin(phi2)), s.E/2)

	m1 := math.Cos(phi1) / math.Sqrt(1-s.E2*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-s.E2*sin2(phi2))

	n := (math.Log(m1) - math.Log(m2)) / (math.Log(t1) - math.Log(t2))
	F := m1 / (n * math.Pow(t1, n))
	rf := s.A * F * math.Pow(tf, n)

	ri := math.Sqrt(math.Pow(east-p.FalseEasting, 2) + math.Pow(rf-(north-p.FalseNorthing), 2))
	if n < 0 && ri > 0 {
		ri = -ri
	}

	ti := math.Pow(ri/(s.A*F), 1/n)

	var theta float64
	if n > 0 {
		theta = math.Atan2((east - p.FalseEasting), (rf - (north - p.FalseNorthing)))
	} else {
		theta = math.Atan2(-(east - p.FalseEasting), -(rf - (north - p.FalseNorthing)))
	}

	phi := math.Pi/2 - 2*math.Atan(ti)

	for i := 0; i < 5; i++ {
		phi = math.Pi/2 - 2*math.Atan(ti*math.Pow((1-s.E*math.Sin(phi))/(1+s.E*math.Sin(phi)), s.E/2))
	}

	lambda := theta/n + lambdaf

	return degree(lambda), degree(phi)
}

func (p LambertConformalConic2SP) FromGeographic(lon, lat float64, s Spheroid) (east, north float64) {
	phi := radian(lat)
	phif := radian(p.LatitudeOfOrigin)
	phi1 := radian(p.StandardParallel1)
	phi2 := radian(p.StandardParallel2)
	lambda := radian(lon)
	lambdaf := radian(p.CentralMeridian)

	t := math.Tan(math.Pi/4-phi/2) / math.Pow((1-s.E*math.Sin(phi))/(1+s.E*math.Sin(phi)), s.E/2)
	tf := math.Tan(math.Pi/4-phif/2) / math.Pow((1-s.E*math.Sin(phif))/(1+s.E*math.Sin(phif)), s.E/2)
	t1 := math.Tan(math.Pi/4-phi1/2) / math.Pow((1-s.E*math.Sin(phi1))/(1+s.E*math.Sin(phi1)), s.E/2)
	t2 := math.Tan(math.Pi/4-phi2/2) / math.Pow((1-s.E*math.Sin(phi2))/(1+s.E*math.Sin(phi2)), s.E/2)

	m1 := math.Cos(phi1) / math.Sqrt(1-s.E2*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-s.E2*sin2(phi2))

	n := (math.Log(m1) - math.Log(m2)) / (math.Log(t1) - math.Log(t2))
	F := m1 / (n * math.Pow(t1, n))
	r := s.A * F * math.Pow(t, n)
	rf := s.A * F * math.Pow(tf, n)
	theta := n * (lambda - lambdaf)

	east = p.FalseEasting + r*math.Sin(theta)
	north = p.FalseNorthing + rf - r*math.Cos(theta)

	return east, north
}

type AlbersConicEqualArea struct {
	LongitudeOfCenter float64
	LatitudeOfCenter  float64
	StandardParallel1 float64
	StandardParallel2 float64
	FalseEasting      float64
	FalseNorthing     float64
}

func (p AlbersConicEqualArea) ToGeographic(east, north float64, s Spheroid) (lon, lat float64) {
	lambdaf := radian(p.LongitudeOfCenter)
	phif := radian(p.LatitudeOfCenter)
	phi1 := radian(p.StandardParallel1)
	phi2 := radian(p.StandardParallel2)

	alphaf := (1 - s.E2) * ((math.Sin(phif) / (1 - s.E2*sin2(phif))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phif))/(1+s.E*math.Sin(phif))))
	alpha1 := (1 - s.E2) * ((math.Sin(phi1) / (1 - s.E2*sin2(phi1))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phi1))/(1+s.E*math.Sin(phi1))))
	alpha2 := (1 - s.E2) * ((math.Sin(phi2) / (1 - s.E2*sin2(phi2))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phi2))/(1+s.E*math.Sin(phi2))))

	m1 := math.Cos(phi1) / math.Sqrt(1-s.E2*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-s.E2*sin2(phi2))

	n := (math.Pow(m1, 2) - math.Pow(m2, 2)) / (alpha2 - alpha1)
	c := math.Pow(m1, 2) + (n * alpha1)
	rf := (s.A * math.Sqrt(c-n*alphaf)) / n

	ri := math.Sqrt(math.Pow(east-p.FalseEasting, 2) + math.Pow(rf-(north-p.FalseNorthing), 2))
	alphai := (c - (math.Pow(ri, 2) * math.Pow(n, 2) / s.A2)) / n
	betai := math.Asin(alphai / (1 - ((1 - s.E2) / (2 * s.E) * math.Log((1-s.E)/(1+s.E)))))

	var theta float64

	if n > 0 {
		theta = math.Atan2((east - p.FalseEasting), (rf - (north - p.FalseNorthing)))
	} else {
		theta = math.Atan2(-(east - p.FalseEasting), -(rf - (north - p.FalseNorthing)))
	}

	phi := betai + ((s.E2/3 + 31*s.E4/180 + 517*s.E6/5040) * math.Sin(2*betai)) + ((23*s.E4/360 + 251*s.E6/3780) * math.Sin(4*betai)) + ((761 * s.E6 / 45360) * math.Sin(6*betai))
	lambda := lambdaf + (theta / n)

	return degree(lambda), degree(phi)
}

func (p AlbersConicEqualArea) FromGeographic(lon, lat float64, s Spheroid) (east, north float64) {
	lambda := radian(lon)
	lambdaf := radian(p.LongitudeOfCenter)
	phi := radian(lat)
	phif := radian(p.LatitudeOfCenter)
	phi1 := radian(p.StandardParallel1)
	phi2 := radian(p.StandardParallel2)

	alpha := (1 - s.E2) * ((math.Sin(phi) / (1 - s.E2*sin2(phi))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phi))/(1+s.E*math.Sin(phi))))
	alphaf := (1 - s.E2) * ((math.Sin(phif) / (1 - s.E2*sin2(phif))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phif))/(1+s.E*math.Sin(phif))))
	alpha1 := (1 - s.E2) * ((math.Sin(phi1) / (1 - s.E2*sin2(phi1))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phi1))/(1+s.E*math.Sin(phi1))))
	alpha2 := (1 - s.E2) * ((math.Sin(phi2) / (1 - s.E2*sin2(phi2))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phi2))/(1+s.E*math.Sin(phi2))))

	m1 := math.Cos(phi1) / math.Sqrt(1-s.E2*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-s.E2*sin2(phi2))

	n := (math.Pow(m1, 2) - math.Pow(m2, 2)) / (alpha2 - alpha1)

	c := math.Pow(m1, 2) + (n * alpha1)

	theta := n * (lambda - lambdaf)
	r := (s.A * math.Sqrt(c-n*alpha)) / n
	rf := (s.A * math.Sqrt(c-n*alphaf)) / n

	east = p.FalseEasting + r*math.Sin(theta)
	north = p.FalseNorthing + rf - r*math.Cos(theta)

	return east, north
}

type LambertAzimuthalEqualArea struct {
	LatitudeOfCenter  float64
	LongitudeOfCenter float64
	FalseEasting      float64
	FalseNorthing     float64
}

func (p LambertAzimuthalEqualArea) ToGeographic(east, north float64, s Spheroid) (lon, lat float64) {
	phi0 := radian(p.LatitudeOfCenter)
	lambda0 := radian(p.LongitudeOfCenter)

	q0 := (1 - s.E2) * ((math.Sin(phi0) / (1 - s.E2*sin2(phi0))) - (1 / (2 * s.E) * math.Log((1-s.E*math.Sin(phi0))/(1+s.E*math.Sin(phi0)))))
	qp := (1 - s.E2) * ((1 / (1 - s.E2)) - ((1 / (2 * s.E)) * math.Log((1-s.E)/(1+s.E))))

	beta0 := math.Asin(q0 / qp)
	rq := s.A * math.Sqrt(qp/2)
	g := s.A * (math.Cos(phi0) / math.Sqrt(1-s.E2*sin2(phi0))) / (rq * math.Cos(beta0))

	rho := math.Sqrt(math.Pow((east-p.FalseEasting)/g, 2) + math.Pow(g*(north-p.FalseNorthing), 2))
	c := 2 * math.Asin(rho/(2*rq))
	betai := math.Asin((math.Cos(c) * math.Sin(beta0)) + ((g * (north - p.FalseNorthing) * math.Sin(c) * math.Cos(beta0)) / rho))

	phi := betai + ((s.E2/3 + (31*s.E4)/180 + (517*s.E6)/5040) * math.Sin(2*betai)) +
		((23*s.E4)/360+(251*s.E6)/3780)*math.Sin(4*betai) +
		((761*s.E6)/45360)*math.Sin(6*betai)
	lambda := lambda0 + math.Atan2((east-p.FalseEasting)*math.Sin(c), (g*rho*math.Cos(beta0)*math.Cos(c)-math.Pow(g, 2)*(north-p.FalseNorthing)*math.Sin(beta0)*math.Sin(c)))

	return degree(lambda), degree(phi)
}

func (p LambertAzimuthalEqualArea) FromGeographic(lon, lat float64, s Spheroid) (east, north float64) {
	phi0 := radian(p.LatitudeOfCenter)
	lambda0 := radian(p.LongitudeOfCenter)

	phi := radian(50.0)
	lambda := radian(5.0)

	q := (1 - s.E2) * ((math.Sin(phi) / (1 - s.E2*sin2(phi))) - (1 / (2 * s.E) * math.Log((1-s.E*math.Sin(phi))/(1+s.E*math.Sin(phi)))))
	q0 := (1 - s.E2) * ((math.Sin(phi0) / (1 - s.E2*sin2(phi0))) - (1 / (2 * s.E) * math.Log((1-s.E*math.Sin(phi0))/(1+s.E*math.Sin(phi0)))))
	qp := (1 - s.E2) * ((1 / (1 - s.E2)) - ((1 / (2 * s.E)) * math.Log((1-s.E)/(1+s.E))))

	beta0 := math.Asin(q0 / qp)
	beta := math.Asin(q / qp)
	rq := s.A * math.Sqrt(qp/2)
	g := s.A * (math.Cos(phi0) / math.Sqrt(1-s.E2*sin2(phi0))) / (rq * math.Cos(beta0))
	b := rq * math.Sqrt(2/(1+math.Sin(beta0)*math.Sin(beta)+(math.Cos(beta0)*math.Cos(beta)*math.Cos(lambda-lambda0))))

	east = p.FalseEasting + ((b * g) * (math.Cos(beta) * math.Sin(lambda-lambda0)))
	north = p.FalseNorthing + (b/g)*((math.Cos(beta0)*math.Sin(beta))-(math.Sin(beta0)*math.Cos(beta)*math.Cos(lambda-lambda0)))

	return east, north
}

type Krovak struct {
	LatitudeOfCenter       float64
	LongitudeOfCenter      float64
	Azimuth                float64
	PseudoStandardParallel float64
	ScaleFactor            float64
	FalseEasting           float64
	FalseNorthing          float64
}

func (p Krovak) FromGeographic(lon, lat float64, s Spheroid) (east, north float64) {
	phic := radian(p.LatitudeOfCenter)
	lambda0 := radian(p.LongitudeOfCenter)
	phi := radian(lat)
	phip := radian(p.PseudoStandardParallel)
	lambda := radian(lon)
	alphac := radian(p.Azimuth)

	A := s.A * math.Sqrt(1-s.E2) / (1 - s.E2*sin2(phic))
	B := math.Sqrt(1 + (s.E2 * math.Pow(math.Cos(phic), 4) / (1 - s.E2)))
	gamma0 := math.Asin(math.Sin(phic) / B)
	t0 := math.Tan(math.Pi/4+gamma0/2) * math.Pow((1+s.E*math.Sin(phic))/(1-s.E*math.Sin(phic)), s.E*B/2) / math.Pow(math.Tan(math.Pi/4+phic/2), B)
	n := math.Sin(phip)
	r0 := p.ScaleFactor * A / math.Tan(phip)

	U := 2 * (math.Atan(t0*math.Pow(math.Tan(phi/2+math.Pi/4), B)/math.Pow((1+s.E*math.Sin(phi))/(1-s.E*math.Sin(phi)), s.E*B/2)) - math.Pi/4)
	V := B * (lambda0 - lambda)
	T := math.Asin(math.Cos(alphac)*math.Sin(U) + math.Sin(alphac)*math.Cos(U)*math.Cos(V))
	D := math.Asin(math.Cos(U) * math.Sin(V) / math.Cos(T))
	theta := n * D
	r := r0 * math.Pow(math.Tan(math.Pi/4+phip/2), n) / math.Pow(math.Tan(T/2+math.Pi/4), n)
	Xp := r * math.Cos(theta)
	Yp := r * math.Sin(theta)

	return -(Yp + p.FalseEasting), -(Xp + p.FalseNorthing)
}

func (p Krovak) ToGeographic(east, north float64, s Spheroid) (lon, lat float64) {
	phic := radian(p.LatitudeOfCenter)
	phip := radian(p.PseudoStandardParallel)
	lambda0 := radian(p.LongitudeOfCenter)
	alphac := radian(p.Azimuth)

	A := s.A * math.Sqrt(1-s.E2) / (1 - s.E2*sin2(phic))
	B := math.Sqrt(1 + (s.E2 * math.Pow(math.Cos(phic), 4) / (1 - s.E2)))
	gamma0 := math.Asin(math.Sin(phic) / B)
	t0 := math.Tan(math.Pi/4+gamma0/2) * math.Pow((1+s.E*math.Sin(phic))/(1-s.E*math.Sin(phic)), s.E*B/2) / math.Pow(math.Tan(math.Pi/4+phic/2), B)
	n := math.Sin(phip)
	r0 := p.ScaleFactor * A / math.Tan(phip)

	Xpi := (-north) - p.FalseNorthing
	Ypi := (-east) - p.FalseEasting
	ri := math.Sqrt(math.Pow(Xpi, 2) + math.Pow(Ypi, 2))
	thetai := math.Atan2(Ypi, Xpi)
	Di := thetai / math.Sin(phip)
	Ti := 2 * (math.Atan(math.Pow(r0/ri, 1/n)*math.Tan(math.Pi/4+phip/2)) - math.Pi/4)
	Ui := math.Asin(math.Cos(alphac)*math.Sin(Ti) - math.Sin(alphac)*math.Cos(Ti)*math.Cos(Di))
	Vi := math.Asin(math.Cos(Ti) * math.Sin(Di) / math.Cos(Ui))

	phi := Ui

	for i := 0; i < 3; i++ {
		phi = 2 * (math.Atan(math.Pow(t0, -1/B)*math.Pow(math.Tan(Ui/2+math.Pi/4), 1/B)*math.Pow((1+s.E*math.Sin(phi))/(1-s.E*math.Sin(phi)), s.E/2)) - math.Pi/4)
	}

	lambda := lambda0 - Vi/B

	return degree(lambda), degree(phi)
}

type Spheroid struct {
	A, Fi                                          float64
	A2, F, F2, B, E2, E, E4, E6, Ei, Ei2, Ei3, Ei4 float64
}

func (s Spheroid) ToGeocentric(lon, lat, h float64) (x, y, z float64) {
	n := s.A / math.Sqrt(1-s.E2*math.Pow(math.Sin(radian(lat)), 2))

	x = (n + h) * math.Cos(radian(lon)) * math.Cos(radian(lat))
	y = (n + h) * math.Cos(radian(lat)) * math.Sin(radian(lon))
	z = (n*math.Pow(s.A*(1-s.F), 2)/(s.A2) + h) * math.Sin(radian(lat))

	return x, y, z
}

func (s Spheroid) FromGeocentric(x, y, z float64) (lon, lat, h float64) {
	sd := math.Sqrt(x*x + y*y)
	T := math.Atan(z * s.A / (sd * s.B))
	B := math.Atan((z + s.E2*(s.A2)/s.B*
		math.Pow(math.Sin(T), 3)) / (sd - s.E2*s.A*math.Pow(math.Cos(T), 3)))
	n := s.A / math.Sqrt(1-s.E2*math.Pow(math.Sin(B), 2))
	h = sd/math.Cos(B) - n
	lon = degree(math.Atan2(y, x))
	lat = degree(B)

	return lon, lat, h
}

func NewSpheroid(a, fi float64) Spheroid {
	s := Spheroid{
		A:  a,
		Fi: fi,
	}

	s.A2 = s.A * s.A
	s.F = 1 / s.Fi
	s.F2 = s.F * s.F
	s.B = s.A * (1 - s.F)
	s.E2 = 2/s.Fi - s.F2
	s.E = math.Sqrt(s.E2)
	s.E4 = s.E2 * s.E2
	s.E6 = s.E4 * s.E2
	s.Ei = (1 - math.Sqrt(1-s.E2)) / (1 + math.Sqrt(1-s.E2))
	s.Ei2 = s.Ei * s.Ei
	s.Ei3 = s.Ei2 * s.Ei
	s.Ei4 = s.Ei3 * s.Ei

	return s
}

func sin2(r float64) float64 {
	return math.Pow(math.Sin(r), 2)
}

func degree(r float64) float64 {
	return r * 180 / math.Pi
}

func radian(g float64) float64 {
	return g * math.Pi / 180
}
